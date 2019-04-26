defmodule DES do
  @moduledoc """
  Encrypt and Decrypt files using DES symmetric algorithm.

  ## Examples
  
  ``` bash
  $ ./bin enc my_file my_file.enc
  $ ./bin dec my_file.enc my_file
  ```

  iex> DES.encrypt(['12345678', '12345678'])
  iex> DES.decrypt(['12345678', '12345678'])

  """

  @block_size_bytes 8
  @num_rounds 16
 
  # @key_size_bytes 7

  @cp_1 [
    57, 49, 41, 33, 25, 17, 9,
    1, 58, 50, 42, 34, 26, 18,
    10, 2, 59, 51, 43, 35, 27,
    19, 11, 3, 60, 52, 44, 36,
    63, 55, 47, 39, 31, 23, 15,
    7, 62, 54, 46, 38, 30, 22,
    14, 6, 61, 53, 45, 37, 29,
    21, 13, 5, 28, 20, 12, 4
  ]

  @e [
    32, 1, 2, 3, 4, 5,
    4, 5, 6, 7, 8, 9,
    8, 9, 10, 11, 12, 13,
    12, 13, 14, 15, 16, 17,
    16, 17, 18, 19, 20, 21,
    20, 21, 22, 23, 24, 25,
    24, 25, 26, 27, 28, 29,
    28, 29, 30, 31, 32, 1
  ]

  @pi [
    58, 50, 42, 34, 26, 18, 10, 2,
    60, 52, 44, 36, 28, 20, 12, 4, 
    62, 54, 46, 38, 30, 22, 14, 6, 
    64, 56, 48, 40, 32, 24, 16, 8, 
    57, 49, 41, 33, 25, 17, 9, 1, 
    59, 51, 43, 35, 27, 19, 11, 3, 
    61, 53, 45, 37, 29, 21, 13, 5, 
    63, 55, 47, 39, 31, 23, 15, 7
  ]

  @pi_1 [
    40, 8, 48, 16, 56, 24, 64, 32,
    39, 7, 47, 15, 55, 23, 63, 31,
    38, 6, 46, 14, 54, 22, 62, 30,
    37, 5, 45, 13, 53, 21, 61, 29,
    36, 4, 44, 12, 52, 20, 60, 28,
    35, 3, 43, 11, 51, 19, 59, 27,
    34, 2, 42, 10, 50, 18, 58, 26,
    33, 1, 41, 9, 49, 17, 57, 25
  ]

  @shift_order [1, 1, 2, 2, 2, 2, 2, 2, 1, 2, 2, 2, 2, 2, 2, 1]

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

  defp generate_keys(key) do
    k = permute(key, @cp_1)
    {left, right} = split_block(k)

    {left, right}
  end

  defp expansion(key) do
    permute(key, @e)
  end

  defp to_binary_list_string(block) do
    Enum.map(block, fn x -> Integer.to_string(x, 2) |> String.pad_leading(@block_size_bytes,"0") end)
  end

  defp to_binary_string(block) do
    to_binary_list_string(block) |> Enum.join |> String.graphemes
  end

  defp permute(block, table) do
    binary_string = to_binary_string block
    permuted = Enum.map(table, fn x -> Enum.at(binary_string, x - 1) end)

    Enum.map(Enum.chunk_every(permuted, @block_size_bytes), fn x -> String.to_integer(Enum.join(x),2) end)
  end

  defp initial_permutation(block) do
    permute block, @pi
  end

  defp final_permutation(block) do
    permute block, @pi_1
  end

  # Read a file, given a path, show content on stdout and return it
  defp read_file(file_path) do
    text = File.read!file_path |> String.trim

    IO.puts "============================== Text ==========================="
    IO.puts text
    IO.puts "==============================================================="

    text
  end

  # Write a content to a file
  defp write_to(file_path, content) do
    File.write file_path, content
    IO.puts "Result saved to file '#{file_path}'..."
  end

  defp split_block(block) do
    block_items = (to_charlist block) |> List.to_string
    left = String.slice(block_items, 0..(div(String.length(block_items),2)-1)) |> to_charlist
    right = String.slice(block_items, (div(String.length(block_items),2))..-1) |> to_charlist

    {left, right}
  end

  # Substitution stage of the algorithm
  defp substitute(left, right) do
    for {x, y} <- (Enum.zip left, round_function(right, 42)), do: x ^^^ y 
  end

  defp round_function(right, key) do
    Enum.map(right, fn x -> x + key end)
  end

  defp encrypt_block(block, n \\ 0)

  defp encrypt_block(block, n) when n == 0 do
    block = initial_permutation block

    {left, right} = split_block(block)
    left = substitute left, right
    block = right ++ left 

    encrypt_block(block, n + 1)
  end

  defp encrypt_block(block, n) when n == @num_rounds do
    final_permutation block
  end

  defp encrypt_block(block, n) do
    {left, right} = split_block(block)
    left = substitute left, right
    block = right ++ left 

    encrypt_block(block, n + 1)
  end

  # Encrypt the given plain text, which must be a an list of blocks
  defp encrypt_blocks(plain) do
    Enum.map(plain, &encrypt_block/1)
  end

  defp decrypt_block(block, n \\ 0)

  defp decrypt_block(block, n) when n == 0 do
    block =  initial_permutation block

    {left, right} = split_block(block)
    left = substitute left, right
    block = right ++ left

    decrypt_block(block, n + 1)
  end

  defp decrypt_block(block, n) when n == @num_rounds do
    final_permutation block
  end

  defp decrypt_block(block, n) do
    {left, right} = split_block(block)
    left = substitute left, right
    block = right ++ left

    decrypt_block(block, n + 1)
  end

  defp split_blocks(text) do
    leftover = List.duplicate (hd ' '), @block_size_bytes
    Enum.chunk_every text, @block_size_bytes, @block_size_bytes, leftover
  end

  # Descrypt the given plain text, which must be a list of blocks
  defp decrypt_blocks(cyphered) do
    Enum.map(cyphered, &decrypt_block/1)
  end

  def encrypt(plain) do
    IO.puts "Encrypting..."
    split_blocks(plain) |> encrypt_blocks |> Enum.join
  end

  def decrypt(plain) do
    IO.puts "Decrypting..."
    split_blocks(plain) |> decrypt_blocks |> Enum.join
  end

  @doc """
  Used for calling the module as an stand-alone file
  """
  def main() do
    {operation, in_file, out_file} = parse_input()
    text = read_file(in_file) |> String.trim |> to_charlist

    process_function = fn
      text when operation == :enc ->
        encrypt text
      text ->
        decrypt text
    end

    write_to(out_file, process_function.(text))
  end
end

DES.main()
