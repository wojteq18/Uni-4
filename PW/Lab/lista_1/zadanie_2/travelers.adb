with Ada.Text_IO;                    use Ada.Text_IO;
with Ada.Numerics.Float_Random;      use Ada.Numerics.Float_Random;
with Random_Seeds;                   use Random_Seeds;
with Ada.Real_Time;                  use Ada.Real_Time;
with Ada.Characters.Handling;        use Ada.Characters.Handling; -- For To_Lower

procedure Travelers_Synchronized is

   -- Travelers moving on the board
   Nr_Of_Travelers : constant Integer := 15;

   Min_Steps : constant Integer := 10;
   Max_Steps : constant Integer := 100;

   Min_Delay : constant Duration := 0.01;
   Max_Delay : constant Duration := 0.05;

   -- 2D Board with torus topology
   Board_Width  : constant Integer := 15; -- m
   Board_Height : constant Integer := 15; -- n

   -- Timing
   Start_Time : Time := Clock;  -- global starting time

   -- Random seeds for the tasks' random number generators
   Seeds : Seed_Array_Type(1 .. Nr_Of_Travelers) := Make_Seeds(Nr_Of_Travelers);

   -- Types, procedures and functions

   -- Positions on the board
   type Position_Type is record
      X : Integer range 0 .. Board_Width - 1;
      Y : Integer range 0 .. Board_Height - 1;
   end record;

   -- Protected Object for Squares
   protected type Protected_Square is
      entry Try_Acquire; -- Blocks if square is occupied
      procedure Release;
   private
      Is_Occupied : Boolean := False;
   end Protected_Square;

   protected body Protected_Square is
      entry Try_Acquire when not Is_Occupied is
      begin
         Is_Occupied := True;
      end Try_Acquire;

      procedure Release is
      begin
         Is_Occupied := False;
      end Release;
   end Protected_Square;

   -- The Board
   type Board_Type is array (0 .. Board_Width - 1, 0 .. Board_Height - 1)
     of Protected_Square;
   Board : Board_Type; -- Global board instance

   -- elementary steps (Unchanged from original)
   procedure Move_Down(Position : in out Position_Type) is
   begin
      Position.Y := (Position.Y + 1) mod Board_Height;
   end Move_Down;
   procedure Move_Up(Position : in out Position_Type) is
   begin
      Position.Y := (Position.Y + Board_Height - 1) mod Board_Height;
   end Move_Up;
   procedure Move_Right(Position : in out Position_Type) is
   begin
      Position.X := (Position.X + 1) mod Board_Width;
   end Move_Right;
   procedure Move_Left(Position : in out Position_Type) is
   begin
      Position.X := (Position.X + Board_Width - 1) mod Board_Width;
   end Move_Left;

   -- traces of travelers (Unchanged from original)
   type Trace_Type is record
      Time_Stamp : Duration;
      Id         : Integer;
      Position   : Position_Type;
      Symbol     : Character;
   end record;

   type Trace_Array_type is array (0 .. Max_Steps) of Trace_Type;

   type Traces_Sequence_Type is record
      Last        : Integer := -1;
      Trace_Array : Trace_Array_type;
   end record;

   procedure Print_Trace(Trace : Trace_Type) is
   begin
      Put_Line(
         Duration'Image(Trace.Time_Stamp) & " " &
           Integer'Image(Trace.Id) & " " &
           Integer'Image(Trace.Position.X) & " " &
           Integer'Image(Trace.Position.Y) & " " &
           (1 => Trace.Symbol) -- Use aggregate to print single char as string
        );
   end Print_Trace;

   procedure Print_Traces(Traces : Traces_Sequence_Type) is
   begin
      for I in 0 .. Traces.Last loop
         Print_Trace(Traces.Trace_Array(I));
      end loop;
   end Print_Traces;

   -- task Printer collects and prints reports of traces (Modified)
   task Printer is
      entry Report(Traces : Traces_Sequence_Type);
   end Printer;

   task body Printer is
      Reports_Received : Natural := 0;
   begin
      loop
         select
            accept Report(Traces : Traces_Sequence_Type) do
               Print_Traces(Traces);
               Reports_Received := Reports_Received + 1;
            end Report;

            if Reports_Received = Nr_Of_Travelers then
               exit; -- Exit loop once all reports are received
            end if;
         or
            terminate; -- Allow termination when all travelers are done
         end select;
      end loop;
   end Printer;

   -- travelers (Unchanged record type)
   type Traveler_Type is record
      Id       : Integer;
      Symbol   : Character;
      Position : Position_Type;
   end record;

   -- traveler task type (Modified body)
   task type Traveler_Task_Type is
      entry Init(Id : Integer; Seed : Integer; Symbol : Character);
      entry Start;
   end Traveler_Task_Type;

   task body Traveler_Task_Type is
      G                : Generator;
      Traveler         : Traveler_Type;
      Nr_of_Steps      : Integer;
      Traces           : Traces_Sequence_Type;
      Current_Position : Position_Type; -- Store current position explicitly
      Stuck            : Boolean := False;

      -- Timeout slightly longer than max random delay to suspect deadlock
      Timeout_Duration : constant Duration := Max_Delay * 2.0;

      procedure Store_Trace is
         Time_Stamp : constant Duration := To_Duration(Clock - Start_Time);
      begin
         Traces.Last := Traces.Last + 1;
         if Traces.Last > Max_Steps then -- Prevent array index error
            Put_Line("Error: Traveler " & Integer'Image(Traveler.Id) & " exceeded Max_Steps for trace storage.");
            Stuck := True; -- Treat as stuck to report and terminate
            return;
         end if;
         Traces.Trace_Array(Traces.Last) := (
           Time_Stamp => Time_Stamp,
           Id         => Traveler.Id,
           Position   => Traveler.Position,
           Symbol     => Traveler.Symbol
           );
      end Store_Trace;

      function Choose_Direction return Integer is
      begin
         return Integer(Float'Floor(4.0 * Random(G)));
      end Choose_Direction;

      function Get_Next_Position(Current : Position_Type; Direction : Integer)
        return Position_Type
      is
         Next_Pos : Position_Type := Current;
      begin
         case Direction is
            when 0 => Move_Up(Next_Pos);
            when 1 => Move_Down(Next_Pos);
            when 2 => Move_Left(Next_Pos);
            when 3 => Move_Right(Next_Pos);
            when others => null;
         end case;
         return Next_Pos;
      end Get_Next_Position;

      procedure Handle_Stuck is
      begin
         Stuck := True;
         Traveler.Symbol := To_Lower(Traveler.Symbol); -- Use standard function
         Store_Trace; -- Store final trace at the stuck position
         Printer.Report(Traces); -- Send report and finish
      end Handle_Stuck;

   begin -- Task body main execution
      accept Init(Id : Integer; Seed : Integer; Symbol : Character) do
         Reset(G, Seed);
         Traveler.Id     := Id;
         Traveler.Symbol := Symbol;

         -- Find a free starting position (can block if highly contested)
         loop
            Traveler.Position := (
              X => Integer(Float'Floor(Float(Board_Width) * Random(G))),
              Y => Integer(Float'Floor(Float(Board_Height) * Random(G)))
              );
            -- Attempt to acquire the initial square
            select
               Board(Traveler.Position.X, Traveler.Position.Y).Try_Acquire;
               exit; -- Got the square
            or
               delay 0.001; -- Small delay if contested, try another spot
            end select;
         end loop;

         Current_Position := Traveler.Position; -- Store initial position
         Store_Trace;                          -- Store starting position trace

         Nr_of_Steps := Min_Steps + Integer(Float'Floor(Float(Max_Steps - Min_Steps + 1) * Random(G)));

      end Init;

      accept Start do
         null;
      end Start;

      step_loop : for Step in 1 .. Nr_of_Steps loop
         delay Min_Delay + (Max_Delay - Min_Delay) * Duration(Random(G));

         declare
            Direction     : Integer       := Choose_Direction;
            Next_Position : Position_Type := Get_Next_Position(Current_Position, Direction);
         begin
            -- Attempt to acquire the next square with a timeout
            select
               Board(Next_Position.X, Next_Position.Y).Try_Acquire;

               -- Success! Release the old square and move.
               Board(Current_Position.X, Current_Position.Y).Release;
               Traveler.Position := Next_Position;
               Current_Position  := Next_Position;
               Store_Trace;

            or
               delay Timeout_Duration;
               -- Timeout occurred - suspect deadlock or contention
               Handle_Stuck;
               exit step_loop; -- Exit the main loop for this traveler
            end select;
         end; -- declare block
      end loop step_loop;

      -- Send report if not already sent due to being stuck
      if not Stuck then
         Board(Current_Position.X, Current_Position.Y).Release; -- Release final square
         Printer.Report(Traces);
      end if;

   end Traveler_Task_Type;

   -- local for main task
   Travel_Tasks : array (0 .. Nr_Of_Travelers - 1) of Traveler_Task_Type;
   Symbol       : Character := 'A';

begin -- Main procedure execution

   -- Print the header line for the display script:
   Put_Line(
      "-1 " &
        Integer'Image(Nr_Of_Travelers) & " " &
        Integer'Image(Board_Width) & " " &
        Integer'Image(Board_Height)
     );

   -- Init travelers tasks
   for I in Travel_Tasks'Range loop
      Travel_Tasks(I).Init(I, Seeds(I + 1), Symbol); -- Seeds index is 1-based
      if Symbol < 'Z' then
         Symbol := Character'Succ(Symbol);
      else -- Wrap around if needed (more than 26 travelers)
         Symbol := 'A';
      end if;
   end loop;

   -- Start travelers tasks
   for I in Travel_Tasks'Range loop
      Travel_Tasks(I).Start;
   end loop;

   -- Main task implicitly waits here for Printer task to finish,
   -- which waits for all travelers to report (or terminate).

end Travelers_Synchronized;