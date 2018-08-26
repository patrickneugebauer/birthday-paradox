function simulate {
  $start = Get-Date
  $iterations = 1000
  $sample_size = 23
  $count = 0

  For ($i=0; $i -lt $iterations; $i++) {
    $data = @(-1) * 23
    For ($l=0; $l -lt $sample_size; $l++) {
      $number = GET-RANDOM -Minimum 0 -Maximum 364
      If ($data.Contains($number)) {
        $count++
        break
      } ELSE {
        $data[$l] = $number
      }
    }
  }
  Write-Host "iterations: $iterations "
  Write-Host "sample-size: $sample_size"
  $percent = $count / $iterations * 100
  Write-Host "percent: $percent"
  $end = Get-Date
  $diff = ($end - $start).TotalSeconds
  Write-Host "seconds: $([Math]::Round($diff, 3))"
}

simulate
