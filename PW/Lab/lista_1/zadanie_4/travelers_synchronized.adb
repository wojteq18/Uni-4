with Ada.Text_IO;                    use Ada.Text_IO;
with Ada.Numerics.Float_Random;      use Ada.Numerics.Float_Random;
with Random_Seeds;                   use Random_Seeds;
with Ada.Real_Time;                  use Ada.Real_Time;
with Ada.Characters.Handling;        use Ada.Characters.Handling; -- For To_Lower

procedure Travelers_Synchronized is

   Nr_Of_Travelers : constant Integer := 15;
   Board_Width     : constant Integer := 15;
   Board_Height    : constant Integer := 15;

   Min_Steps : constant Integer := 10;
   Max_Steps : constant Integer := 100;

   Min_Delay : constant Duration := 0.01;
   Max_Delay : constant Duration := 0.05;

   Start_Time : Time := Clock;  -- global starting time

   Seeds : Seed_Array_Type(1 .. Nr_Of_Travelers) := Make_Seeds(Nr_Of_Travelers);

   type Position_Type is record
      X : Integer range 0 .. Board_Width - 1;
      Y : Integer range 0 .. Board_Height - 1;
   end record;


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

   type Board_Type is array (0 .. Board_Width - 1, 0 .. Board_Height - 1)
     of Protected_Square;
   Board : Board_Type; -- Global board instance

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
           (1 => Trace.Symbol)
        );
   end Print_Trace;

   procedure Print_Traces(Traces : Traces_Sequence_Type) is
   begin
      for I in 0 .. Traces.Last loop
         Print_Trace(Traces.Trace_Array(I));
      end loop;
   end Print_Traces;

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
               exit; 
            end if;
         or
            terminate; 
         end select;
      end loop;
   end Printer;

   type Traveler_Type is record
      Id       : Integer;
      Symbol   : Character;
      Position : Position_Type;
   end record;

   task type Traveler_Task_Type is
      entry Init(Id : Integer; Seed : Integer; Symbol : Character);
      entry Start;
   end Traveler_Task_Type;

   task body Traveler_Task_Type is
      G                : Generator;
      Traveler         : Traveler_Type;
      Nr_of_Steps      : Integer;
      Traces           : Traces_Sequence_Type;
      Current_Position : Position_Type;
      Stuck            : Boolean := False;

      Fixed_Direction  : Integer := 0; 

      Timeout_Duration : constant Duration := Max_Delay * 2.0;

      procedure Store_Trace is
         Time_Stamp : constant Duration := To_Duration(Clock - Start_Time);
      begin
         Traces.Last := Traces.Last + 1;
         if Traces.Last > Max_Steps then
            Put_Line("Error: Traveler " & Integer'Image(Traveler.Id) &
                     " exceeded Max_Steps for trace storage.");
            Stuck := True;
            return;
         end if;
         Traces.Trace_Array(Traces.Last) := (
           Time_Stamp => Time_Stamp,
           Id         => Traveler.Id,
           Position   => Traveler.Position,
           Symbol     => Traveler.Symbol
         );
      end Store_Trace;

      procedure Handle_Stuck is
      begin
         Stuck := True;
         Traveler.Symbol := To_Lower(Traveler.Symbol);
         Store_Trace; -- final trace
         Printer.Report(Traces);
      end Handle_Stuck;

      function Get_Next_Position
        (Current : Position_Type;
         Direction : Integer) return Position_Type
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

   begin -- Traveler task body
      accept Init(Id : Integer; Seed : Integer; Symbol : Character) do
         Reset(G, Seed);
         Traveler.Id     := Id;
         Traveler.Symbol := Symbol;

         Traveler.Position := (X => Id, Y => Id);

         Board(Traveler.Position.X, Traveler.Position.Y).Try_Acquire;

         Current_Position := Traveler.Position;
         Store_Trace;

         Nr_of_Steps :=
           Min_Steps
           + Integer(
               Float'Floor(
                 Float(Max_Steps - Min_Steps + 1) * Random(G)
               )
             );

         if (Traveler.Id mod 2) = 0 then  -- parzysty
            if Random(G) < 0.5 then
               Fixed_Direction := 0; -- up
            else
               Fixed_Direction := 1; -- down
            end if;
         else
            if Random(G) < 0.5 then
               Fixed_Direction := 2; -- left
            else
               Fixed_Direction := 3; -- right
            end if;
         end if;

      end Init;

      accept Start do
         null;
      end Start;

      step_loop : for Step in 1 .. Nr_of_Steps loop
         -- losowe opóźnienie
         delay Min_Delay + (Max_Delay - Min_Delay) * Duration(Random(G));

         declare
            Next_Position : Position_Type :=
              Get_Next_Position(Current_Position, Fixed_Direction);
         begin
            -- Spróbuj przejąć kwadrat docelowy z timeoutem
            select
               Board(Next_Position.X, Next_Position.Y).Try_Acquire;
               -- Zwolnij starą pozycję
               Board(Current_Position.X, Current_Position.Y).Release;
               Traveler.Position := Next_Position;
               Current_Position  := Next_Position;
               Store_Trace;
            or
               delay Timeout_Duration;
               -- Uznajemy, że deadlock/timeout
               Handle_Stuck;
               exit step_loop;
            end select;
         end;
      end loop step_loop;

      if not Stuck then
         -- Zwolnij ostateczny square
         Board(Current_Position.X, Current_Position.Y).Release;
         Printer.Report(Traces);
      end if;
   end Traveler_Task_Type;

   Travel_Tasks : array (0 .. Nr_Of_Travelers - 1) of Traveler_Task_Type;
   Symbol       : Character := 'A';

begin
   Put_Line("-1 " &
            Integer'Image(Nr_Of_Travelers) & " " &
            Integer'Image(Board_Width) & " " &
            Integer'Image(Board_Height));

   for I in Travel_Tasks'Range loop
      Travel_Tasks(I).Init(I, Seeds(I + 1), Symbol);
      if Symbol < 'Z' then
         Symbol := Character'Succ(Symbol);
      else
         Symbol := 'A';
      end if;
   end loop;

   for I in Travel_Tasks'Range loop
      Travel_Tasks(I).Start;
   end loop;

end Travelers_Synchronized;
