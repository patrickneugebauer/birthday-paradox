import 'dart:math';

void main(List<String> arguments) {
  var iterations = int.parse(arguments[0]);
  simulate(iterations);
}

void simulate(int iterations) {
  final start = new DateTime.now().millisecondsSinceEpoch;
  const sampleSize = 23;

  var count = 0;
  final random = new Random();
  for (var i=0; i < iterations; i++) {
    final arr = new List.filled(365, 0);
    for (var j=0; j < sampleSize; j++) {
      final number = random.nextInt(365);
      if (arr[number] == 1) {
        count++;
        break;
      } else {
        arr[number] = 1;
      }
    }
  }
  print('iterations: $iterations');
  print('sample-size: $sampleSize');
  final results = (count / iterations * 100).toStringAsFixed(2);
  print('percent: $results');
  final end = new DateTime.now().millisecondsSinceEpoch;
  final diff = ((end - start) / 1000).toStringAsFixed(3);
  print('seconds: $diff');
}
