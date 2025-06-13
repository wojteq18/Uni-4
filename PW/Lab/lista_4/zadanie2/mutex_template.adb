with Ada.Text_IO; use Ada.Text_IO;
with Ada.Numerics.Float_Random;
with Random_Seeds; use Random_Seeds;
with Ada.Real_Time; use Ada.Real_Time;
with Ada.Strings.Unbounded; use Ada.Strings.Unbounded;
with Ada.Integer_Text_IO; use Ada.Integer_Text_IO;

procedure Mutex_Szymanski is -- Zmieniono nazwę procedury, aby pasowała do pliku

  -- Processes
  Nr_Of_Processes : constant Integer := 15;

  Min_Steps : constant Integer := 50;
  Max_Steps : constant Integer := 100;

  Min_Delay : constant Duration := 0.01;
  Max_Delay : constant Duration := 0.05;
  Small_Yield_Delay : constant Duration := 0.001;

  type Process_State is (
    Local_Section,
    Entry_Protocol_1,
    Entry_Protocol_2,
    Entry_Protocol_3,
    Entry_Protocol_4,
    Critical_Section,
    Exit_Protocol
    );

  Board_Width  : constant Integer := Nr_Of_Processes;
  Board_Height : constant Integer := Process_State'Pos(Process_State'Last) + 1;

  Start_Time : Time := Clock;

  Seeds : Seed_Array_Type(1 .. Nr_Of_Processes) := Make_Seeds(Nr_Of_Processes);

  -- << TUTAJ DEKLARACJA NOWEGO TYPU TABLICOWEGO >>
  type Flags_Array_Internal_Type is array (0 .. Nr_Of_Processes - 1) of Integer range 0 .. 4;

  protected Flags_Manager is
    procedure Set_Flag(Process_Index : Integer; Value : Integer);
    function Get_Flag(Process_Index : Integer) return Integer;
  private
    -- << UŻYCIE NAZWANEGO TYPU >>
    Flags_Array : Flags_Array_Internal_Type := (others => 0);
  end Flags_Manager;

  protected body Flags_Manager is
    procedure Set_Flag(Process_Index : Integer; Value : Integer) is
    begin
      Flags_Array(Process_Index) := Value;
    end Set_Flag;

    function Get_Flag(Process_Index : Integer) return Integer is
    begin
      return Flags_Array(Process_Index);
    end Get_Flag;
  end Flags_Manager;


  -- Types, procedures and functions
  type Position_Type is record
    X : Integer range 0 .. Board_Width - 1;
    Y : Integer range 0 .. Board_Height - 1;
  end record;

  type Trace_Type is record
    Time_Stamp : Duration;
    Id         : Integer;
    Position   : Position_Type;
    Symbol     : Character;
  end record;

  type Trace_Array_Type is array (0 .. Max_Steps) of Trace_Type; -- Increased size slightly just in case

  type Traces_Sequence_Type is record
    Last        : Integer := -1;
    Trace_Array : Trace_Array_Type;
  end record;


  procedure Print_Trace(Trace : Trace_Type) is
  begin
    Put_Line(
      Duration'Image(Trace.Time_Stamp) & " " &
      Integer'Image(Trace.Id) & " " &
      Integer'Image(Trace.Position.X) & " " &
      Integer'Image(Trace.Position.Y) & " " &
      (1 => Trace.Symbol)  -- Print as a 1-character string
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
  begin
    for I in 1 .. Nr_Of_Processes loop
      accept Report(Traces : Traces_Sequence_Type) do
        Print_Traces(Traces);
      end Report;
    end loop;

    Put(
      "-1 " &
      Integer'Image(Nr_Of_Processes) & " " &
      Integer'Image(Board_Width) & " " &
      Integer'Image(Board_Height) & " "
      );
    for I in Process_State'Range loop
      Put(Process_State'Image(I) & ";"); -- Use 'Image for enum
    end loop;
    Put_Line("EXTRA_LABEL;");
  end Printer;


  type Process_Info_Type is record -- Renamed from Process_Type to avoid conflict
    Id       : Integer;
    Symbol   : Character;
    Position : Position_Type;
  end record;


  task type Process_Task_Type is
    entry Init(Id : Integer; Seed_Val : Integer; Symb : Character); -- Renamed params
    entry Start;
  end Process_Task_Type;

  task body Process_Task_Type is
    G            : Ada.Numerics.Float_Random.Generator; -- Explicit package
    Process_Data : Process_Info_Type; -- Use renamed type
    Current_Time_Stamp : Duration; -- Renamed from Time_Stamp
    Nr_Of_Steps  : Integer;
    Traces       : Traces_Sequence_Type;
    Current_State : Process_State := Local_Section; -- Track current state for Change_State

    procedure Store_Trace is
    begin
      if Traces.Last < Traces.Trace_Array'Last then
        Traces.Last := Traces.Last + 1;
        Traces.Trace_Array(Traces.Last) := (
          Time_Stamp => Current_Time_Stamp,
          Id         => Process_Data.Id,
          Position   => Process_Data.Position,
          Symbol     => Process_Data.Symbol);
      else
         Ada.Text_IO.Put_Line("Error: Trace_Array full for process " & Integer'Image(Process_Data.Id));
      end if;
    end Store_Trace;

    procedure Change_State_And_Trace(New_State : Process_State) is
    begin
      Current_State := New_State;
      Current_Time_Stamp := To_Duration(Clock - Start_Time);
      Process_Data.Position.Y := Process_State'Pos(Current_State);
      Store_Trace;
    end Change_State_And_Trace;


  begin
    accept Init(Id : Integer; Seed_Val : Integer; Symb : Character) do
      Ada.Numerics.Float_Random.Reset(G, Seed_Val); -- Explicit package
      Process_Data.Id := Id;
      Process_Data.Symbol := Symb;
      Process_Data.Position := (
        X => Id,
        Y => Process_State'Pos(Local_Section));
      Nr_Of_Steps := Min_Steps + Integer(Float(Max_Steps - Min_Steps) * Ada.Numerics.Float_Random.Random(G));
      Current_Time_Stamp := To_Duration(Clock - Start_Time);
      Store_Trace; -- store starting position
    end Init;

    accept Start do
      null;
    end Start;

    -- Consistent loop iterations with Go example (approx. 7 stages per CS entry)
    for Step in 0 .. (Nr_Of_Steps / 7) - 1 loop

      -- LOCAL_SECTION
      Change_State_And_Trace(Local_Section); -- Ensure state is traced
      delay Min_Delay + (Max_Delay - Min_Delay) * Duration(Ada.Numerics.Float_Random.Random(G));

      -- ENTRY_PROTOCOL - Szymański's Algorithm
      -- Step 1: Flag[id] := 1
      Flags_Manager.Set_Flag(Process_Data.Id, 1);
      Change_State_And_Trace(Entry_Protocol_1);

      -- Step 2: await (∀k: Flag[k] < 3)
      loop -- Await loop
        declare
          Must_Still_Wait_S2 : Boolean := False;
        begin
          for J in 0 .. Nr_Of_Processes - 1 loop
            if Flags_Manager.Get_Flag(J) = 3 or Flags_Manager.Get_Flag(J) = 4 then
              Must_Still_Wait_S2 := True;
              exit; -- Exit inner loop
            end if;
          end loop;
          exit when not Must_Still_Wait_S2; -- Exit await loop if condition met
          delay Small_Yield_Delay;
        end;
      end loop; -- End await loop for Step 2

      -- Step 3: Flag[id] := 3
      Flags_Manager.Set_Flag(Process_Data.Id, 3);
      Change_State_And_Trace(Entry_Protocol_3);

      -- Step 4: if (∃k ≠ id: Flag[k] = 1) then { Flag[id] := 2; await (Flag[k] ≠ 4); }
      declare
        Found_With_Flag_1_S4 : Boolean := False;
        Idx_K_S4             : Integer := -1;
      begin
        for K_Loop_Var in 0 .. Nr_Of_Processes - 1 loop
          if K_Loop_Var /= Process_Data.Id then
            if Flags_Manager.Get_Flag(K_Loop_Var) = 1 then
              Found_With_Flag_1_S4 := True;
              Idx_K_S4 := K_Loop_Var;
              exit;
            end if;
          end if;
        end loop;

        if Found_With_Flag_1_S4 then
          Flags_Manager.Set_Flag(Process_Data.Id, 2);
          Change_State_And_Trace(Entry_Protocol_2);
          -- Ada.Text_IO.Put_Line("Process " & Process_Data.Id'Image & " was in Entry_Protocol_2!!!"); -- Debug

          loop -- Await Flag(Idx_K_S4) /= 4
            exit when Flags_Manager.Get_Flag(Idx_K_S4) /= 4;
            delay Small_Yield_Delay;
          end loop;
        end if;
      end; -- End declare block for Step 4

      -- Step 5: Flag[id] := 4 (unconditionally after step 4 logic)
      Flags_Manager.Set_Flag(Process_Data.Id, 4);
      Change_State_And_Trace(Entry_Protocol_4);

      -- Step 6: await (∀k < id: Flag[k] < 2)
      loop -- Await loop
        declare
          Must_Still_Wait_S6 : Boolean := False;
        begin
          for J in 0 .. Process_Data.Id - 1 loop
            if Flags_Manager.Get_Flag(J) >= 2 then -- Wait if 2, 3, or 4
              Must_Still_Wait_S6 := True;
              exit;
            end if;
          end loop;
          exit when not Must_Still_Wait_S6;
          delay Small_Yield_Delay;
        end;
      end loop; -- End await loop for Step 6

      -- Step 7: await (∀k > id: Flag[k] < 2 ∨ Flag[k] > 3) (i.e., wait if Flag[k] = 2 or Flag[k] = 3)
      loop -- Await loop
        declare
          Must_Still_Wait_S7 : Boolean := False;
        begin
          for K_Loop_Var in Process_Data.Id + 1 .. Nr_Of_Processes - 1 loop
            declare
              Fk : constant Integer := Flags_Manager.Get_Flag(K_Loop_Var);
            begin
              if Fk = 2 or Fk = 3 then
                Must_Still_Wait_S7 := True;
                exit;
              end if;
            end;
          end loop;
          exit when not Must_Still_Wait_S7;
          delay Small_Yield_Delay; -- Using Small_Yield_Delay consistently
        end;
      end loop; -- End await loop for Step 7

      -- CRITICAL_SECTION
      Change_State_And_Trace(Critical_Section);
      delay Min_Delay + (Max_Delay - Min_Delay) * Duration(Ada.Numerics.Float_Random.Random(G));

      -- EXIT_PROTOCOL
      Change_State_And_Trace(Exit_Protocol);
      Flags_Manager.Set_Flag(Process_Data.Id, 0); -- Szymański's exit is just setting flag to 0

      -- Back to LOCAL_SECTION (implicitly handled by Change_State_And_Trace at start of next loop iteration)
      -- Explicitly trace the return to Local_Section for clarity if preferred, or ensure it's traced at loop start
      -- Current code will trace Local_Section at the beginning of the next iteration.
      -- To match Go version's exact end-of-loop tracing:
      Change_State_And_Trace(Local_Section);

    end loop;

    Printer.Report(Traces);
  end Process_Task_Type;


  Process_Tasks : array (0 .. Nr_Of_Processes - 1) of Process_Task_Type;
  Symbol_Counter : Character := 'A'; -- Renamed from Symbol to avoid conflict

begin
  for I in Process_Tasks'Range loop
    Process_Tasks(I).Init(Id => I, Seed_Val => Seeds(I + 1), Symb => Symbol_Counter);
    if Symbol_Counter < Character'Last then
       Symbol_Counter := Character'Succ(Symbol_Counter);
    end if;
  end loop;

  for I in Process_Tasks'Range loop
    Process_Tasks(I).Start;
  end loop;

end Mutex_Szymanski;
