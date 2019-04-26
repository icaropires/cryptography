defmodule DesTest do
  use ExUnit.Case
  doctest Des

  test "greets the world" do
    assert Des.hello() == :world
  end
end
