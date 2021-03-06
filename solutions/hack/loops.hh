<?hh

function simulate($iterations) {
  $start = microtime(true);
  $sample_size = 23;

  $count = 0;
  for ($x = 0; $x < $iterations; $x++) {
    $data = array_fill(0, 365, 0);
    for ($n = 0; $n < $sample_size; $n++) {
      $number = rand(0, 364);
      if ($data[$number] === 1) {
        $count++;
        break;
      } else {
        $data[$number] = 1;
      }
    }
  }
  print "iterations: $iterations\n";
  print "sample-size: $sample_size\n";
  $percent = round($count / $iterations * 100, 2);
  print "percent: $percent\n";
  $end = microtime(true);
  $seconds = round($end - $start, 3);
  print "seconds: $seconds\n";
}

<<__EntryPoint>>
function main() {
  // $iterations = $_GLOBALS['argv'][1];
  $iterations = HH\global_get('argv')[1];
  simulate($iterations);
}
