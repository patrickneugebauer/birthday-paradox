import 'dart:math';

void main() {
  simulate();
}

void simulate() {
  final start = new DateTime.now().millisecondsSinceEpoch;
  const iterations = 1000 * 1000;
  const sampleSize = 23;

  var count = 0;
  final random = new Random();
  for (var i=0; i < iterations; i++) {
    final arr = [];
    for (var j=0; j < sampleSize; j++) {
      final number = random.nextInt(364);
      if (arr.contains(number)) {
        count++;
        break;
      } else {
        arr.add(number);
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
