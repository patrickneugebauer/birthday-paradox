program Loops;
uses sysutils, math, dateutils;

const
  SAMPLE_SIZE = 23;
var
  start: TDateTime;
  iterations: LongInt;
  // data
  data: array [1..365] of Integer;
  num: Integer;
  i: LongInt;
  s: Integer;
  count: LongInt = 0;
  // calcs
  percent: Single;
  fin: TDateTime;
  milliseconds: Comp;
  seconds: Single;
begin
  // data
  start := time;
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
  fin := time;
  milliseconds := (TimeStampToMSecs(DateTimeToTimeStamp(fin)) - TimeStampToMSecs(DateTimeToTimeStamp(start)));
  seconds := milliseconds / 1000;

  // output
  writeln('iterations: ', iterations);
  writeln('sample-size: ', SAMPLE_SIZE);
  writeln('percent: ', FormatFloat('0.00', percent));
  writeln('seconds: ', FormatFloat('0.000', seconds));
end.
