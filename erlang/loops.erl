-module(loops).
-export([main/1]).

-define(SampleSize, 23).

% functions
random_day(_) -> rand:uniform(365).
random_sample(N) -> lists:map(fun(X) -> random_day(X) end, lists:seq(1, N)).
% identity(X) -> X. % cannot get (fun identity/1) to work
% ---------------------------------------------------------------------------
% performance notes
% recursive implementation is 4.5x slower than set implementation in escript
% but 50% faster when compiled
% stated another way: the recursive implementation gets >50x faster when compiled
% but the set implementation gets <10x faster
has_duplicates([]) -> false;
has_duplicates([X|XS]) -> case lists:member(X, XS) of true -> true; false -> has_duplicates(XS) end.
% has_duplicates(List) -> erlang:length(List) == sets:size(sets:from_list(List)).
% ---------------------------------------------------------------------------

main(Args) ->
  % data
  StartMicrosecond = os:system_time(),
  Iterations = list_to_integer(lists:nth(1, Args)),
  Data = lists:map(
    fun(_) -> has_duplicates(random_sample(?SampleSize)) end,
    lists:seq(1, Iterations)
  ),
  Duplicates = erlang:length(lists:filter(fun(X) -> X end, Data)),
  Percent = Duplicates / Iterations * 100,
  EndMicrosecond = os:system_time(),
  Seconds = (EndMicrosecond - StartMicrosecond) / 1000 / 1000 / 1000,
  % output
  io:fwrite("iterations: ~w~n", [Iterations]),
  io:fwrite("sample-size: ~w~n", [?SampleSize]),
  io:fwrite("percent: ~.2f~n", [Percent]),
  io:fwrite("seconds: ~.3f~n", [Seconds]).
