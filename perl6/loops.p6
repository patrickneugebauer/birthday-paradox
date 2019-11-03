#!/usr/bin/perl

sub MAIN($iterations) {
  my $start = now;
  my $sample_size = 23;
  my $count = 0;

  for 1..$iterations {
    my @data = (for 1..365 { 0 });
    sample: for 1..$sample_size {
      my $number = Int(365.rand);
      if (@data[$number] == 1) {
        $count++;
        last sample;
      } else {
        @data[$number] = 1;
      }
    }
  }

  say "iterations: $iterations";
  say "sample-size: $sample_size";
  my $percent = $count / $iterations * 100;
  printf("percent: %.2f\n", $percent);
  my $end = now;
  my $diff = $end - $start;
  printf("seconds: %.3f\n", $diff);
}
