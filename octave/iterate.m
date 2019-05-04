function iterate
  start = time();
  iterations = 10000;
  sample_size = 23;
  day_count = 365;
  count = 0;
  matrix = randi(365, iterations, sample_size);
  for i = 1:iterations
    data = matrix(i, :);
    if (columns(unique(data)) != sample_size)
      count++;
    endif
  endfor
  printf("iterations: %d\n", iterations)
  printf("sample-size: %d\n", sample_size)
  percent = count / iterations * 100;
  printf("percent: %.2f\n", percent)
  finish = time();
  diff = finish - start;
  printf("seconds: %.3f\n", diff)
endfunction
