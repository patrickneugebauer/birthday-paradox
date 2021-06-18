defmodule Elixir.MixProject do
  use Mix.Project

  def project do
    [
      app: :loops,
      version: "0.1.0",
      escript: escript
    ]
  end

  defp escript do
    [main_module: Loops.CLI]
  end
end
