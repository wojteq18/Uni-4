with Ada.Text_IO; use Ada.Text_IO;
with Ada.Numerics.Float_Random;
with Random_Seeds; use Random_Seeds;
with Ada.Real_Time;
with Monitor_Package; -- Nasz poprawiony pakiet

procedure Readers_Writers_Demo is

   Min_Steps : constant Integer := 5;  -- Zmniejszone dla szybszych testów
   Max_Steps : constant Integer := 10;
   Min_Delay : constant Duration := 0.01;
   Max_Delay : constant Duration := 0.05;

   -- Definicje Czytelników i Pisarzy
   Nr_Of_Readers   : constant Integer := 10;
   Nr_Of_Writers   : constant Integer := 5;
   Nr_Of_Processes : constant Integer := Nr_Of_Readers + Nr_Of_Writers;

   -- Stany Procesu
   type Process_Role is (Reader, Writer);
   type Process_State is (
      Local_Section,
      Start,          -- Wchodzenie do czytelni
      Reading_Room,   -- W czytelni
      Stop            -- Wychodzenie z czytelni
   );

   Board_Width  : constant Integer := Nr_Of_Processes;
   Board_Height : constant Integer := Process_State'Pos( Process_State'Last ) + 1;

   Start_Time : Ada.Real_Time.Time := Ada.Real_Time.Clock;
   Seeds : Seed_Array_Type( 1..Nr_Of_Processes ) := Make_Seeds( Nr_Of_Processes );

   type Position_Type is record
      X: Integer range 0 .. Board_Width - 1;
      Y: Integer range 0 .. Board_Height - 1;
   end record;

   type Trace_Type is record
      Time_Stamp:  Duration;
      Id : Integer;
      Position: Position_Type;
      Symbol: Character;
   end record;

   type Trace_Array_type is array(0 .. Max_Steps * 4) of Trace_Type; -- *4 dla stanów

   type Traces_Sequence_Type is record
      Last: Integer := -1;
      Trace_Array: Trace_Array_type ;
   end record;

   procedure Print_Trace( Trace : Trace_Type ) is
   begin
      Put_Line(
          Duration'Image( Trace.Time_Stamp ) & " " &
          Integer'Image( Trace.Id ) & " " &
          Integer'Image( Trace.Position.X ) & " " &
          Integer'Image( Trace.Position.Y ) & " " &
          (1 => Trace.Symbol)
        );
   end Print_Trace;

   procedure Print_Traces( Traces : Traces_Sequence_Type ) is
   begin
      for I in 0 .. Traces.Last loop
         Print_Trace( Traces.Trace_Array( I ) );
      end loop;
   end Print_Traces;

   task Printer is
      entry Report( Traces : Traces_Sequence_Type );
   end Printer;

   task body Printer is
   begin
      for I in 1 .. Nr_Of_Processes loop
         accept Report( Traces : Traces_Sequence_Type ) do
            Print_Traces( Traces );
         end Report;
      end loop;

      Put(
         "-1 "&
         Integer'Image( Nr_Of_Processes ) &" "&
         Integer'Image( Board_Width ) &" "&
         Integer'Image( Board_Height ) &" "
      );
      for I in Process_State'Range loop
         Put( Process_State'Image( I ) & ";" );
      end loop;
      Put_Line("EXTRA_LABEL;");
   end Printer;

   -- Zmienne współdzielone dla logiki Czytelników/Pisarzy
   Active_Readers   : Natural := 0;
   Active_Writers   : Natural := 0;
   Waiting_Writers  : Natural := 0; -- Dla priorytetu pisarzy

   OK_To_Read  : aliased Monitor_Package.Condition; -- Aliased jeśli przekazujemy jako 'in out'
   OK_To_Write : aliased Monitor_Package.Condition; -- do procedury Wait

   -- Procedury monitora dla Czytelników/Pisarzy
   procedure Start_Read is
   begin
      Monitor_Package.Monitor.Enter;
      -- Pisarze mają priorytet: jeśli pisarz pisze LUB czeka, czytelnik czeka.
      while Active_Writers > 0 or Waiting_Writers > 0 loop
         Monitor_Package.Wait(OK_To_Read); -- Używamy opakowanej procedury Wait
      end loop;
      Active_Readers := Active_Readers + 1;
      -- Kaskadowe budzenie innych czytelników, jeśli to bezpieczne
      if Monitor_Package.Non_Empty(OK_To_Read) and Active_Writers = 0 and Waiting_Writers = 0 then
          Monitor_Package.Signal(OK_To_Read);
      end if;
      Monitor_Package.Monitor.Leave;
   end Start_Read;

   procedure Stop_Read is
   begin
      Monitor_Package.Monitor.Enter;
      Active_Readers := Active_Readers - 1;
      if Active_Readers = 0 and Waiting_Writers > 0 then
         Monitor_Package.Signal(OK_To_Write); -- Obudź pisarza, jeśli jest ostatnim czytelnikiem
      end if;
      Monitor_Package.Monitor.Leave;
   end Stop_Read;

   procedure Start_Write is
   begin
      Monitor_Package.Monitor.Enter;
      Waiting_Writers := Waiting_Writers + 1;
      while Active_Readers > 0 or Active_Writers > 0 loop
         Monitor_Package.Wait(OK_To_Write);
      end loop;
      Waiting_Writers := Waiting_Writers - 1;
      Active_Writers := Active_Writers + 1; -- Powinno być 1, bo tylko jeden pisarz naraz
      Monitor_Package.Monitor.Leave;
   end Start_Write;

   procedure Stop_Write is
   begin
      Monitor_Package.Monitor.Enter;
      Active_Writers := Active_Writers - 1; -- Powinno być 0
      -- Kogo obudzić? Daj priorytet innym czekającym pisarzom.
      if Waiting_Writers > 0 then
         Monitor_Package.Signal(OK_To_Write);
      else
         -- Jeśli nie ma pisarzy, obudź czytelników (kaskadowo)
         -- Signal budzi tylko jednego. Potrzebna pętla lub kaskada.
         -- Proste: zasygnalizuj raz, pierwszy obudzony czytelnik w Start_Read zasygnalizuje następnego.
         if Monitor_Package.Non_Empty(OK_To_Read) then
            Monitor_Package.Signal(OK_To_Read);
         end if;
      end if;
      Monitor_Package.Monitor.Leave;
   end Stop_Write;

   -- Task Czytelnika/Pisarza
   task type Process_Task_Type is
      entry Init(Id_Num: Integer; Seed: Integer; Role_Symbol: Character; P_Role : Process_Role);
      entry Start_Execution;
   end Process_Task_Type;

   task body Process_Task_Type is
   G : Ada.Numerics.Float_Random.Generator;
   
   -- Najpierw definiujemy typ rekordu
   type Process_Details_Type is record
      Id       : Integer;
      Symbol   : Character;
      Role     : Process_Role;
      Position : Position_Type;
   end record;
   
   -- Potem deklarujemy zmienną tego typu
   Process_Details : Process_Details_Type;
   
   Time_Stamp : Duration;
   Nr_Of_Local_Loops : Integer;
   Traces: Traces_Sequence_Type;

      procedure Store_Trace is
      begin
         Traces.Last := Traces.Last + 1;
         if Traces.Last < Trace_Array_type'Length then
             Traces.Trace_Array( Traces.Last ) := (
                 Time_Stamp => Time_Stamp,
                 Id => Process_Details.Id,
                 Position => Process_Details.Position,
                 Symbol => Process_Details.Symbol
               );
         else
            Ada.Text_IO.Put_Line("Trace buffer overflow for task " & Integer'Image(Process_Details.Id));
         end if;
      end Store_Trace;

      procedure Change_State( State: Process_State ) is
      begin
         Time_Stamp := Ada.Real_Time.To_Duration ( Ada.Real_Time.Clock - Start_Time );
         Process_Details.Position.Y := Process_State'Pos( State );
         Store_Trace;
      end Change_State;

   begin
      accept Init(Id_Num: Integer; Seed: Integer; Role_Symbol: Character; P_Role : Process_Role) do
         Ada.Numerics.Float_Random.Reset(G, Seed);
         Process_Details.Id := Id_Num;
         Process_Details.Symbol := Role_Symbol;
         Process_Details.Role   := P_Role;
         Process_Details.Position := (
            X => Id_Num -1, -- Indeksowanie od 0 dla pozycji X
            Y => Process_State'Pos( LOCAL_SECTION )
           );
         Nr_Of_Local_Loops := Min_Steps + Integer( Float(Max_Steps - Min_Steps) * Ada.Numerics.Float_Random.Random(G));
         Time_Stamp := Ada.Real_Time.To_Duration ( Ada.Real_Time.Clock - Start_Time );
         Store_Trace;
      end Init;

      accept Start_Execution do
         null;
      end Start_Execution;

      for I in 1 .. Nr_Of_Local_Loops loop
         Change_State( LOCAL_SECTION );
         delay Min_Delay + (Max_Delay - Min_Delay) * Duration(Ada.Numerics.Float_Random.Random(G));

         Change_State( START );
         if Process_Details.Role = Reader then
            Start_Read;
         else -- Writer
            Start_Write;
         end if;

         Change_State( READING_ROOM );
         delay Min_Delay + (Max_Delay - Min_Delay) * Duration(Ada.Numerics.Float_Random.Random(G));

         Change_State( STOP );
         if Process_Details.Role = Reader then
            Stop_Read;
         else -- Writer
            Stop_Write;
         end if;
      end loop;

      Printer.Report( Traces );
   end Process_Task_Type;

   Process_Tasks: array (1 .. Nr_Of_Processes) of Process_Task_Type;

begin
   -- Inicjalizacja tasków
   for I in 1 .. Nr_Of_Readers loop
      Process_Tasks(I).Init( I, Seeds(I), 'R', Reader );
   end loop;
   for I in 1 .. Nr_Of_Writers loop
      Process_Tasks(Nr_Of_Readers + I).Init( Nr_Of_Readers + I, Seeds(Nr_Of_Readers + I), 'W', Writer );
   end loop;

   -- Start tasków
   for I in Process_Tasks'Range loop
      Process_Tasks(I).Start_Execution;
   end loop;

   -- Ciało głównej procedury czeka na zakończenie Printera, który czeka na wszystkie Process_Tasks
   -- To zapewnia, że program nie zakończy się przedwcześnie.
   -- Taski Monitor i Condition zakończą się dzięki klauzuli 'terminate',
   -- gdy wszystkie taski klienckie (Process_Tasks) zakończą działanie,
   -- a Printer zakończy się po zebraniu wszystkich raportów.

end Readers_Writers_Demo;
