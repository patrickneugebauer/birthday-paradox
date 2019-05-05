import System.Random
import Data.Time.Clock.POSIX
import qualified Data.IntSet as IntSet

main = do
  start <- getTimeMillis
  let iterations = 100000 :: Int
  gen <- getStdGen
  let results = (/ (fromIntegral iterations)) . fromIntegral . simulate iterations $ gen
  putStrLn . ("iterations: " ++) . show $ iterations
  putStrLn . ("sample-size: " ++) . show $ 23
  putStrLn . ("percent: " ++) . show . roundTo 2 . (*100) $ results
  end <- getTimeMillis
  putStrLn . ("seconds: " ++) . show . (/1000) . fromIntegral $ (end - start)


getTimeMillis :: (Integral a) => IO a
getTimeMillis = round <$> (*1000) <$> getPOSIXTime

simulate :: Int -> StdGen -> Int
simulate 0 gen = 0
simulate n gen =
  let range = (0, 364) :: (Int, Int)
      sampleSize = 23 :: Int
      sample = IntSet.empty :: IntSet.IntSet
      (x, newGen) = runIteration range sampleSize sample gen
  in x + simulate (n - 1) newGen

runIteration :: (Int, Int) -> Int -> IntSet.IntSet -> StdGen -> (Int, StdGen)
runIteration rng 0 xs gen = (0, gen)
runIteration rng n xs gen =
  let (x, newGen) = randomR rng gen :: (Int, StdGen)
  in if IntSet.member x xs then (1, newGen) else runIteration rng (n - 1) (IntSet.insert x xs) newGen

roundTo n =
  let shifter = 10^n
  in (/shifter) . fromIntegral . round . (*shifter)
