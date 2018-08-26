simulate <- function() {
  start <- Sys.time()
  iterations <- 10000
  sample_size <- 23
  count <- 0
  for (i in 1:iterations) {
    data <- vector(,365)
    for (n in 1:sample_size) {
      number <- sample(1:365,1)
      if (data[number] == 1) {
        count <- count + 1
        break
      } else {
        data[number] <- 1
      }
    }
  }
  cat(paste("iterations:", iterations, "\n"))
  cat(paste("sample-size:", sample_size, "\n"))
  percent <- round(count / iterations * 100, digits = 2)
  cat(paste("percent:", percent, "\n"))
  end <- Sys.time()
  diff <- round(as.numeric(end - start), digits = 3)
  cat(paste("seconds:", diff, "\n"))
}

simulate()
