function iterate
  start = time();
  iterations = 500;
  sample_size = 23;
  day_count = 365;
  count = 0;
  for i = 1:iterations
    % need to create a matrix in each iteration
    % octave uses non-zero indexed arrays with paren lookup
    % data = []
    data = zeros(1, day_count);
    for s = 1:sample_size
      num = randi(day_count);
      if (data(num) == 1)
        count++;
        break;
      else
        data(num) = 1;
      endif
    endfor
  endfor
  disp(cstrcat("iterations: ", mat2str(iterations)))
  disp(cstrcat("sample-size: ", mat2str(sample_size)))
  percent = count / iterations * 100;
  disp(cstrcat("percent: ", mat2str(percent, 2)))
  finish = time();
  diff = finish - start;
  disp(cstrcat("seconds: ", mat2str(diff, 3)))
endfunction
