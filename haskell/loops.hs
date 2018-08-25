import System.Random
import Data.Time.Clock.POSIX

getTimeMillis :: (Integral a) => IO a
getTimeMillis = round <$> (*1000) <$> getPOSIXTime

duplicates :: (Eq a) => [a] -> Bool
duplicates [] = False
duplicates (x:xs) = x `elem` xs || duplicates xs

roundTo :: (Integral a, RealFrac b, Fractional c) => a -> b -> c
roundTo n = (/10^n) . fromIntegral . round . (*10^n)

randomRs' :: (RandomGen g, Random a, Eq a, Num n, Eq n) => (a, a) -> n -> g -> ([a], g)
randomRs' _ 0 gen = ([], gen)
randomRs' (l, u) n gen =
  let (r, midGen) = randomR (l, u) gen
      (xs, finalGen) = randomRs' (l, u) (n-1) midGen
  in (r:xs, finalGen)

simulate :: (Num n, Eq n) => n -> StdGen -> [Bool]
simulate 0 gen = []
simulate n gen =
  let (rs, midGen) = randomRs' (0, 364) 23 gen :: ([Int], StdGen)
  in duplicates rs:simulate (n-1) midGen

simulate' :: Int -> StdGen -> Int
simulate' 0 gen = 0
simulate' n gen =
  let range = (0, 364) :: (Int, Int)
      sampleSize = 23 :: Int
      sample = [] :: [Int]
      (x, newGen) = runIteration range sampleSize sample gen
  in x + simulate' (n - 1) newGen

runIteration :: (Int, Int) -> Int -> [Int] -> StdGen -> (Int, StdGen)
runIteration rng 0 xs gen = (0, gen)
runIteration rng n xs gen =
  let (x, newGen) = randomR rng gen :: (Int, StdGen)
  in if x `elem` xs then (1, newGen) else runIteration rng (n - 1) (x:xs) newGen

main = do
  start <- getTimeMillis
  let iterations = 100000 :: Int
  gen <- getStdGen
  -- let results = (/ (fromIntegral iterations)) . fromIntegral . length . filter id . simulate iterations $ gen
  let results = (/ (fromIntegral iterations)) . fromIntegral . simulate' iterations $ gen
  putStrLn . ("iterations: " ++) . show $ iterations
  putStrLn . ("sample-size: " ++) . show $ 23
  putStrLn . ("percent: " ++) . show . roundTo 2 . (*100) $ results
  end <- getTimeMillis
  putStrLn . ("seconds: " ++) . show . (/1000) . fromIntegral $ (end - start)
