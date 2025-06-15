-- Copyright (C) 2006 M. Ben-Ari. See copyright.txt
-- Modified by AI based on user's request to address potential race condition.
-- Date: 2025-06-15
-- User: wojteq18

package body Monitor_Package is

  task body Monitor is
  begin
    loop
      -- Monitor po prostu zapewnia wzajemne wykluczanie
      accept Enter;
      accept Leave;
    end loop;
  end Monitor;

  task body Condition is
    -- Lokalny licznik "rezerwacji" i faktycznie oczekujących w Wait.
    -- Zastępuje bezpośrednie użycie Wait'Count w strażnikach, aby uniknąć race condition.
    My_Count : Natural := 0;
  begin
    loop
      select
        -- Akceptuj Pre_Wait: proces sygnalizuje zamiar czekania
        accept Pre_Wait do
          My_Count := My_Count + 1;
        end Pre_Wait;
      or
        -- Akceptuj Signal: tylko jeśli ktoś faktycznie czeka lub zarezerwował (My_Count > 0)
        when My_Count > 0 =>
          accept Signal do
            -- Proces sygnalizujący musi opuścić monitor, aby pozwolić
            -- obudzonemu procesowi ponownie wejść.
            -- My_Count zostanie zdekrementowany, gdy proces opuści Wait.
            -- LUB, jeśli Signal jest akceptowany w pętli wewnętrznej Wait,
            -- to My_Count jest dekrementowany tam. Zgodnie z poleceniem,
            -- dekrementacja w Signal w pętli wewnętrznej Wait.
            null; -- Samo 'Signal' nie zwalnia monitora bezpośrednio, robi to procedura opakowująca
                  -- lub kod klienta. Monitor.Leave jest teraz odpowiedzialnością kodu klienta
                  -- lub procedury opakowującej Wait.
                  -- W oryginalnej implementacji Signal zawierało Monitor.Leave.
                  -- Tutaj, aby być spójnym z nową logiką Wait, Signal tylko budzi.
                  -- Zgodnie z poleceniem, Signal w pętli wewnętrznej Wait dekrementuje My_Count.
          end Signal;
      or
        -- Akceptuj Wait: tylko jeśli proces wcześniej wywołał Pre_Wait (My_Count > 0 sugeruje to)
        -- Strażnik 'when My_Count > 0' jest tu ważny, aby Wait nie blokowało się bez Pre_Wait.
        when My_Count > 0 =>
          accept Wait do
            -- Proces jest teraz formalnie w stanie oczekiwania wewnątrz 'Wait'.
            -- Pętla wewnętrzna do obsługi sygnału lub sprawdzenia stanu.
            loop
              select
                accept Signal do
                  -- Proces został obudzony przez Signal.
                  My_Count := My_Count - 1; -- Zmniejszamy licznik, bo proces opuszcza Wait.
                  exit; -- Wyjdź z wewnętrznej pętli oczekiwania
                end Signal;
              or
                -- Akceptuj Pre_Wait również tutaj, na wypadek gdyby inny proces
                -- próbował zarezerwować, gdy ten jest w pętli wewnętrznej.
                accept Pre_Wait do
                  My_Count := My_Count + 1;
                end Pre_Wait;
              or
                -- Akceptuj Is_Waiting również tutaj
                accept Is_Waiting(B: out Boolean) do
                  B := My_Count > 0; -- Ktoś czeka, jeśli licznik jest dodatni
                end Is_Waiting;
                -- Nie ma potrzeby 'Waiting' z oryginalnej implementacji,
                -- bo Is_Waiting pełni tę rolę.
              end select;
            end loop;
          end Wait;
      or
        -- Akceptuj Is_Waiting: sprawdza, czy ktoś zarezerwował miejsce lub czeka
        accept Is_Waiting(B: out Boolean) do
          B := My_Count > 0;
        end Is_Waiting;
      or
        -- Dodajemy terminate, aby task mógł się zakończyć, jeśli nie jest już potrzebny
        -- (np. gdy obiekt Condition wychodzi poza zakres).
         terminate; -- Opcjonalne, zależy od wymagań zarządzania cyklem życia tasków.
        -- Dla klasycznego monitora, taski żyją tak długo jak program.
        null; -- Aby uniknąć potencjalnego busy-wait, jeśli żaden entry nie jest gotowy.
              -- Lepszym rozwiązaniem byłoby użycie delay, jeśli terminate nie jest opcją.
              -- Ale w typowym scenariuszu monitora, zawsze będzie jakieś wywołanie.
      end select;
    end loop;
  end Condition;

  -- Implementacja funkcji pomocniczej
  function Non_Empty(C: in out Condition) return Boolean is
    B: Boolean;
  begin
    C.Is_Waiting(B); -- Używa nowego entry
    return B;
  end Non_Empty;

  -- Implementacja procedury opakowującej
  procedure Wait(C: in out Condition) is
  begin
    C.Pre_Wait;       -- Krok 1: Zarezerwuj miejsce i zwiększ My_Count
    Monitor.Leave;    -- Krok 2: Opuść monitor
    C.Wait;           -- Krok 3: Faktycznie czekaj (teraz My_Count chroni przed race condition)
                      -- Po powrocie z C.Wait, proces powinien ponownie wejść do monitora,
                      -- co zwykle robi się przez Monitor.Enter na początku procedury monitorowej.
                      -- Ta procedura Wait sama w sobie nie robi Monitor.Enter po C.Wait.
  end Wait;

end Monitor_Package;