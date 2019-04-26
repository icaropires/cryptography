defmodule Euclides do
  def euclides(a, b) when b <= 0 do
    a
  end

  def euclides(a, b) do
    euclides b, (rem a, b)
  end

  # ax + by = d = gcd(a, b)
  def euclides_extended(a, b) when a == 0 do
    {b, 0, 1}
  end

  def euclides_extended(a, b) do
    {gcd, x, y} = euclides_extended((rem b, a), a)
    {gcd, (y - (div b, a) * x), x}
  end

  def print_euclides_extended(a, b) do
    cond do
      b == 0 ->
        IO.puts "Can't divide by zero!!!"
      a == 0 ->
        IO.puts "Zero doesn't have a multiplicativa inverse!!!"
      true ->
        {gcd, x, y} = euclides_extended a, b
        inverse = rem ((rem x, b) + b), b
        IO.puts "GCD = #{gcd}, A = #{a}, B = #{b}, X = #{x}, Y = #{y}, Inverse = #{inverse}"
    end
  end

  def main() do
    # [a, b] = (IO.read :stdio, :line) |> String.trim |> (String.split " ") |> (Enum.map &String.to_integer/1)
    # result = euclides a, b
    
    print_euclides_extended(3041, 17331)
    print_euclides_extended(213, 21753)
    print_euclides_extended(548, 9571)
    print_euclides_extended(24573, 68432)
    print_euclides_extended(10, 0)
    print_euclides_extended(0, 10)
  end
end

Euclides.main()
