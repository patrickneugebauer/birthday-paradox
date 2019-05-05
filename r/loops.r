simulate <- function() {
  start <- Sys.time()
  iterations <- 75000
  sample_size <- 23
  count <- 0
  for (i in 1:iterations) {
    list <- sample(1:365, sample_size, replace=TRUE)
    if (length(unique(list)) != sample_size) count <- count + 1
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
