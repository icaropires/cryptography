defmodule DES do
  @moduledoc """
  Encrypt and Decrypt files using DES symmetric algorithm.

  ## Examples
  
  ``` bash
  $ ./bin enc my_file my_file.enc
  $ ./bin dec my_file.enc my_file
  ```

  iex> DES.encrypt_blocks(['12345678', '12345678'])
  iex> DES.decrypt_blocks(['12345678', '12345678'])

  """

  @block_size_bytes 8
  # @key_size_bytes 7

  require System
  require IO
  require String
  use Bitwise

  # Parse argv inputs to make sure user called the module right. Exits on fail
  defp parse_input() do
    if length(System.argv) == 3 do
      [operation, in_file, out_file]  = System.argv
      IO.puts "Processando arquivo '#{in_file}'..."

      operation = case operation do
        "enc" ->
          :enc
        "dec" ->
          :dec
        _ ->
          IO.puts "Invalid operation !!!"
          exit(:normal)
          ""
      end
      {operation, in_file, out_file}
    else
      IO.puts "Pass the parameters like this:\n./bin [enc|dec] [in_file] [out_file]"
      exit(:normal)
    end
  end

  # Read a file, given a path, show content on stdout and return it
  defp read_file(file_path) do
    text = File.read!file_path |> String.trim

    IO.puts "============================== Text ==========================="
    IO.puts text
    IO.puts "===================================================================="

    text
  end

  # Write a content to a file
  defp write_to(file_path, content) do
    File.write file_path, content
    IO.puts "Result saved to file '#{file_path}'..."
  end

  defp split_block(block) do
    block_items = (List.to_string(to_charlist(block)))
    left = String.slice(block_items, 0..(div(String.length(block_items),2)-1)) |> to_charlist
    right = String.slice(block_items, (div(String.length(block_items),2))..-1) |> to_charlist
    {left, right}
  end

  @doc """
  Permutation stage of the algorithm
  """
  defp permutate(block) do
    {left, right} = split_block(block)
    right ++ left
  end

  @doc """
  Substitution stage of the algorithm
  """
  defp substitute(left, right) do
    for {x, y} <- (Enum.zip left, round_function(right, 42)), do: x ^^^ y 
  end

  defp round_function(right, key) do
    Enum.map(right, fn x -> x + key end )
  end

  def encrypt_block(block, n \\ 0) do
    IO.puts "Encrypting..."
    if n > 16 do
      permutate block
    else
      {left, right} = split_block(block)
      left = substitute left, right
      block = right ++ left 
      encrypt_block(block, n + 1)
    end
  end

  @doc """
  Encrypt the given plain text, which must be a an list of blocks
  """
  def encrypt_blocks(plain, n \\ 0) do
    Enum.map(plain, &encrypt_block/1)
  end

  def decrypt_block(block, n \\ 0) do
    IO.puts "Decrypting..."
    if n > 16 do
      permutate block
    else
      {left, right} = split_block(block)
      left = substitute left, right
      block = right ++ left
      decrypt_block(block, n + 1)
    end
  end

  @doc """
  Descrypt the given plain text, which must be a list of blocks
  """
  def decrypt_blocks(cyphered, n \\ 0) do
    Enum.map(cyphered, &decrypt_block/1)
  end

  @doc """
  Used for calling the module as an stand-alone file
  """
  def main() do
    {operation, in_file, out_file} = parse_input()

    text = (read_file in_file) |> String.trim("\n") |> to_charlist
    blocks = Enum.chunk_every text, @block_size_bytes, @block_size_bytes, List.duplicate(32, 8) 
    IO.puts(tl blocks)

    process_function = fn
      text when operation == :enc ->
        encrypt_blocks text
      text ->
        decrypt_blocks text
    end

    result = process_function.(blocks) |> Enum.join
    write_to out_file, result
  end
end

DES.main()
