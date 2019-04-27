defmodule Euclides do
  use Bitwise

  @spec sum(integer(), integer()) :: integer()
  def sum(a, b) do
    a ^^^ b
  end

  @spec subtract(integer(), integer()) :: integer()
  def subtract(a, b) do
    a ^^^ b
  end

  @spec multiply(integer(), integer()) :: integer()
  def multiply(a, b) do
    a ^^^ b
  end

  @spec divide(integer(), integer()) :: integer()
  def divide(a, b) do
    a ^^^ b
  end

  def main() do
    IO.puts "========================="
    IO.puts "Select the operation:"
    IO.puts "1. Sum"
    IO.puts "2. Subtraction"
    IO.puts "3. Multiplication"
    IO.puts "4. Division"
    IO.puts "========================="
    option = (IO.read :stdio, :line) |> String.trim |> String.to_integer

    IO.puts "Insert the operands: A B"
    [a, b] = (IO.read :stdio, :line) |> String.trim |> (String.split " ") |> (Enum.map &String.to_integer/1)

    result = case option do
      1 ->
        sum a, b
      2 ->
        subtract a, b
      3 ->
        multiply a, b
      4 ->
        divide a, b
      _ ->
        IO.puts "Invalid option!"
    end

    IO.puts result
  end
end

Euclides.main()
