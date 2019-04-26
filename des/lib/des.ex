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

  defp generate_keys(key) do
    shift_order = [1,1,2,2,2,2,2,2,1,2,2,2,2,2,2,1]
    cp_1 = [57, 49, 41, 33, 25, 17, 9,
        1, 58, 50, 42, 34, 26, 18,
        10, 2, 59, 51, 43, 35, 27,
        19, 11, 3, 60, 52, 44, 36,
        63, 55, 47, 39, 31, 23, 15,
        7, 62, 54, 46, 38, 30, 22,
        14, 6, 61, 53, 45, 37, 29,
        21, 13, 5, 28, 20, 12, 4]
    k = permute(key, cp_1)
    {left, right} = split_block(k)
    
  end

  defp expansion(key) do
    e = [32, 1, 2, 3, 4, 5,
     4, 5, 6, 7, 8, 9,
     8, 9, 10, 11, 12, 13,
     12, 13, 14, 15, 16, 17,
     16, 17, 18, 19, 20, 21,
     20, 21, 22, 23, 24, 25,
     24, 25, 26, 27, 28, 29,
     28, 29, 30, 31, 32, 1]
    r = String.graphemes(List.to_string(Enum.map(key, fn x -> Integer.to_string(x,2) |> String.pad_leading(8,"0") end)))
    res = Enum.map(e, fn x -> Enum.at(r,x - 1) end)
    bla = Enum.map(Enum.chunk_every(res, 8), fn x -> String.to_integer(Enum.join(x),2) end)
    IO.puts(bla)
    bla 
  end

  defp permute(block, table) do
    r = String.graphemes(List.to_string(Enum.map(block, fn x -> Integer.to_string(x,2) |> String.pad_leading(8,"0") end)))
    res = Enum.map(table, fn x -> Enum.at(r,x - 1) end)
    bla = Enum.map(Enum.chunk_every(res, 8), fn x -> String.to_integer(Enum.join(x),2) end)
    IO.puts(bla)
    bla 
  end

  defp initial_permutation(block) do
    pi = [58, 50, 42, 34, 26, 18, 10, 2,
      60, 52, 44, 36, 28, 20, 12, 4, 
      62, 54, 46, 38, 30, 22, 14, 6, 
      64, 56, 48, 40, 32, 24, 16, 8, 
      57, 49, 41, 33, 25, 17, 9, 1, 
      59, 51, 43, 35, 27, 19, 11, 3, 
      61, 53, 45, 37, 29, 21, 13, 5, 
      63, 55, 47, 39, 31, 23, 15, 7]
    r = String.graphemes(List.to_string(Enum.map(block, fn x -> Integer.to_string(x,2) |> String.pad_leading(8,"0") end)))
    res = Enum.map(pi, fn x -> Enum.at(r,x - 1) end)
    bla = Enum.map(Enum.chunk_every(res, 8), fn x -> String.to_integer(Enum.join(x),2) end)
    IO.puts(bla)
    bla 
  end

  defp final_permutation(block) do
    pi_1 = [40, 8, 48, 16, 56, 24, 64, 32,
        39, 7, 47, 15, 55, 23, 63, 31,
        38, 6, 46, 14, 54, 22, 62, 30,
        37, 5, 45, 13, 53, 21, 61, 29,
        36, 4, 44, 12, 52, 20, 60, 28,
        35, 3, 43, 11, 51, 19, 59, 27,
        34, 2, 42, 10, 50, 18, 58, 26,
        33, 1, 41, 9, 49, 17, 57, 25] 
    r = String.graphemes(List.to_string(Enum.map(block, fn x -> Integer.to_string(x,2) |> String.pad_leading(8,"0") end)))
    res = Enum.map(pi_1, fn x -> Enum.at(r,x - 1) end)
    bla = Enum.map(Enum.chunk_every(res, 8), fn x -> String.to_integer(Enum.join(x),2) end)
    IO.puts(bla)
    bla 

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
    cond do
      n == 0 ->  
        block = initial_permutation block
        {left, right} = split_block(block)
        left = substitute left, right
        block = right ++ left 
        encrypt_block(block, n + 1)
      n == 16 -> 
        final_permutation block
      true -> 
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
    b = Enum.map(plain, &final_permutation/1)
    a = Enum.map(b, &initial_permutation/1)
    IO.puts("-----------")
    IO.puts(a)
    IO.puts(b)
    IO.puts("-----------")
    Enum.map(plain, &encrypt_block/1)
  end

  def decrypt_block(block, n \\ 0) do
    IO.puts "Decrypting..."
    cond do
      n == 0 ->  
        block =  initial_permutation block
        {left, right} = split_block(block)
        left = substitute left, right
        block = right ++ left
        decrypt_block(block, n + 1)
      n == 16 -> 
        final_permutation block
      true ->  
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
