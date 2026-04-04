program Loops;
uses BaseUnix, Unix, sysutils, math, dateutils;

const
  SAMPLE_SIZE = 23;
var
  start: timeval;
  iterations: LongInt;
  // data
  data: array [1..365] of Integer;
  num: Integer;
  i: LongInt;
  s: Integer;
  count: LongInt = 0;
  // calcs
  percent: Single;
  fin: timeval;
  milliseconds: Comp;
  seconds: Single;
begin
  // data
  fpGetTimeOfDay(@start, nil);
  iterations := StrToInt(ParamStr(1));
  for i := 1 to iterations do
  begin
    FillChar(data, SizeOf(data), 0);
    for s := 1 to SAMPLE_SIZE do
    begin
      num := floor(random * 365) + 1;
      if data[num] = 1 then
        begin
          Inc(count);
          break;
        end
      else
        data[num] := 1;
    end;
  end;

  // calcs
  percent := count / iterations * 100;
  fpGetTimeOfDay(@fin, nil);
  seconds := (fin.tv_sec - start.tv_sec) + (fin.tv_usec - start.tv_usec) / 1000000.0;

  // output
  writeln('iterations: ', iterations);
  writeln('sample-size: ', SAMPLE_SIZE);
  writeln('percent: ', FormatFloat('0.00', percent));
  writeln('seconds: ', FormatFloat('0.000000', seconds));
end.
