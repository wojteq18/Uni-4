with Ada.Text_IO; use Ada.Text_IO;

procedure Hello_Task is

   --Definicja zadania (taska)
   task type Printer_Task is
      entry Print_Message(Msg : String);
   end Printer_Task;

   --Implementacja zadania
   task body Printer_Task is
   begin
      loop
         accept Print_Message(Msg : String) do
            Put_Line("Message: " & Msg);
         end Print_Message;
      end loop;
   end Printer_Task;  

   P1, P2 : Printer_Task;

begin
   P1.Print_Message("Hello from P1");
   P2.Print_Message("Hello from P2");

   P1.Print_Message("Goodbye from P1");
   P2.Print_Message("Goodbye from P2");
end Hello_Task;   