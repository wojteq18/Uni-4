with Ada.Text_IO;                    use Ada.Text_IO;

procedure Hello is 
   function Square(X : Integer) return Integer is
   begin
      return X * X;
   end Square;

   X : Integer := 5;
   Result : Integer := Square(X);
begin
   if X > 4 then
      Put_Line("X is greater than 4");
   else
      Put_Line("X is not greater than 4");
   end if;

   for I in 1 .. 10 loop
      Put_Line("Squere = " & Integer'Image(Result));
   end loop;
end Hello;

