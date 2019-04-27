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

  @p [
    16, 7, 20, 21, 29, 12, 28, 17,
    1, 15, 23, 26, 5, 18, 31, 10,
    2, 8, 24, 14, 32, 27, 3, 9,
    19, 13, 30, 6, 22, 11, 4, 25
  ]

  @s_box [
    [[14, 4, 13, 1, 2, 15, 11, 8, 3, 10, 6, 12, 5, 9, 0, 7],
     [0, 15, 7, 4, 14, 2, 13, 1, 10, 6, 12, 11, 9, 5, 3, 8],
     [4, 1, 14, 8, 13, 6, 2, 11, 15, 12, 9, 7, 3, 10, 5, 0],
     [15, 12, 8, 2, 4, 9, 1, 7, 5, 11, 3, 14, 10, 0, 6, 13],
    ],

    [[15, 1, 8, 14, 6, 11, 3, 4, 9, 7, 2, 13, 12, 0, 5, 10],
     [3, 13, 4, 7, 15, 2, 8, 14, 12, 0, 1, 10, 6, 9, 11, 5],
     [0, 14, 7, 11, 10, 4, 13, 1, 5, 8, 12, 6, 9, 3, 2, 15],
     [13, 8, 10, 1, 3, 15, 4, 2, 11, 6, 7, 12, 0, 5, 14, 9],
    ],

    [[10, 0, 9, 14, 6, 3, 15, 5, 1, 13, 12, 7, 11, 4, 2, 8],
     [13, 7, 0, 9, 3, 4, 6, 10, 2, 8, 5, 14, 12, 11, 15, 1],
     [13, 6, 4, 9, 8, 15, 3, 0, 11, 1, 2, 12, 5, 10, 14, 7],
     [1, 10, 13, 0, 6, 9, 8, 7, 4, 15, 14, 3, 11, 5, 2, 12],
    ],

    [[7, 13, 14, 3, 0, 6, 9, 10, 1, 2, 8, 5, 11, 12, 4, 15],
     [13, 8, 11, 5, 6, 15, 0, 3, 4, 7, 2, 12, 1, 10, 14, 9],
     [10, 6, 9, 0, 12, 11, 7, 13, 15, 1, 3, 14, 5, 2, 8, 4],
     [3, 15, 0, 6, 10, 1, 13, 8, 9, 4, 5, 11, 12, 7, 2, 14],
    ],

    [[2, 12, 4, 1, 7, 10, 11, 6, 8, 5, 3, 15, 13, 0, 14, 9],
     [14, 11, 2, 12, 4, 7, 13, 1, 5, 0, 15, 10, 3, 9, 8, 6],
     [4, 2, 1, 11, 10, 13, 7, 8, 15, 9, 12, 5, 6, 3, 0, 14],
     [11, 8, 12, 7, 1, 14, 2, 13, 6, 15, 0, 9, 10, 4, 5, 3],
    ],

    [[12, 1, 10, 15, 9, 2, 6, 8, 0, 13, 3, 4, 14, 7, 5, 11],
     [10, 15, 4, 2, 7, 12, 9, 5, 6, 1, 13, 14, 0, 11, 3, 8],
     [9, 14, 15, 5, 2, 8, 12, 3, 7, 0, 4, 10, 1, 13, 11, 6],
     [4, 3, 2, 12, 9, 5, 15, 10, 11, 14, 1, 7, 6, 0, 8, 13],
    ],

    [[4, 11, 2, 14, 15, 0, 8, 13, 3, 12, 9, 7, 5, 10, 6, 1],
     [13, 0, 11, 7, 4, 9, 1, 10, 14, 3, 5, 12, 2, 15, 8, 6],
     [1, 4, 11, 13, 12, 3, 7, 14, 10, 15, 6, 8, 0, 5, 9, 2],
     [6, 11, 13, 8, 1, 4, 10, 7, 9, 5, 0, 15, 14, 2, 3, 12],
    ],

    [[13, 2, 8, 4, 6, 15, 11, 1, 10, 9, 3, 14, 5, 0, 12, 7],
     [1, 15, 13, 8, 10, 3, 7, 4, 12, 5, 6, 11, 0, 14, 9, 2],
     [7, 11, 4, 1, 9, 12, 14, 2, 0, 6, 10, 13, 15, 3, 5, 8],
     [2, 1, 14, 7, 4, 10, 8, 13, 15, 12, 9, 0, 3, 5, 6, 11],
    ]
  ]
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

  @cp_2 [
    14, 17, 11, 24, 1, 5, 3, 28,
    15, 6, 21, 10, 23, 19, 12, 4,
    26, 8, 16, 7, 27, 20, 13, 2,
    41, 52, 31, 37, 47, 55, 30, 40,
    51, 45, 33, 48, 44, 49, 39, 56,
    34, 53, 46, 42, 50, 36, 29, 32
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
  @spec parse_input() :: {atom(), String.t(), String.t()}
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

  @spec expansion([integer()]) :: [integer()]
  defp expansion(block) do
    permute block, @e
  end

  defp shifts(list, left, right, index \\ 0)

  defp shifts(list, left, right, index) when index == (length @shift_order) do 
    list 
  end

  defp shifts(list, left, right, index) do
    left = shift(left, Enum.at(@shift_order, index)) 
    right = shift(right, Enum.at(@shift_order, index)) 
    list = list ++ [Enum.map(@cp_2, fn x -> Enum.at((left++right), x - 1) end)]
    shifts(list, left, right, index + 1)

  end

  @spec generate_keys([integer()]) :: [integer()]
  defp generate_keys(key) do
    k = permute(key, @cp_1, false)
    list = []
    left = Enum.slice(k,0..27)
    right = Enum.slice(k,28..56)
    keys = shifts(list, left, right) 
    IO.puts(keys)
    o = Enum.map(keys, fn x -> Enum.map(Enum.chunk_every(x, @block_size_bytes), fn t -> String.to_integer(Enum.join(t),2) end) end)
    o
  end

  @spec shift([String.t()], integer()) :: [integer()]
  defp shift(block, shift_size) do
    Enum.slice(block,shift_size..27) ++ Enum.slice(block,0..(shift_size-1))
  end

  @spec to_binary_list_string(charlist()) :: [String.t()]
  defp to_binary_list_string(block) do
    Enum.map(block, fn x -> Integer.to_string(x, 2) |> String.pad_leading(@block_size_bytes,"0") end)
  end

  @spec to_binary_string(charlist()) :: [String.t()]
  defp to_binary_string(block) do
    to_binary_list_string(block) |> Enum.join |> String.graphemes
  end

  @spec permute(charlist(), [integer()]) :: [integer()]
  defp permute(block, table, to_integer \\ true) do
    binary_string = to_binary_string block
    permuted = Enum.map(table, fn x -> Enum.at(binary_string, x - 1) end)
    if to_integer do
      Enum.map(Enum.chunk_every(permuted, @block_size_bytes), fn x -> String.to_integer(Enum.join(x),2) end)
    else
      permuted
    end
  end

  defp substitute(block) do 
    binary_string = to_binary_string block
    IO.puts(binary_string)
    blocks = Enum.chunk_every(binary_string, 6)
    #IO.inspect blocks
    edges = for i <- blocks, do: [Enum.at(i, 0), Enum.at(i, 5)] |> Enum.join |> String.to_integer(2)
    middles = for k <- blocks, do: Enum.slice(k, 1..4) |> Enum.join |> String.to_integer(2) 
    a = for index <- 0..7, do: Enum.at(@s_box, index) |> Enum.at(Enum.at(edges, index)) |> Enum.at(Enum.at(middles, index))
    IO.puts("=======================")
    a
  end

  @spec initial_permutation(charlist()) :: [integer()]
  defp initial_permutation(block) do
    permute block, @pi
  end

  @spec final_permutation([integer()]) :: [integer()]
  defp final_permutation(block) do
    permute block, @pi_1
  end

  # Read a file, given a path, show content on stdout and return it
  @spec read_file(String.t()) :: String.t()
  defp read_file(file_path) do
    text = File.read!file_path |> String.trim

    IO.puts "============================== Text ==========================="
    IO.puts text
    IO.puts "==============================================================="

    text
  end

  # Write a content to a file
  @spec write_to(String.t(), String.t()) :: atom()
  defp write_to(file_path, content) do
    File.write file_path, content
    IO.puts "Result saved to file '#{file_path}'..."
  end

  @spec split_block([integer()]) :: {[integer()], [integer()]}
  defp split_block(block) do
    half = div (length block), 2
    [left, right] = Enum.chunk_every block, half

    {left, right}
  end

  # Substitution stage of the algorithm
  @spec xor([integer()], [integer()]) :: [integer()]
  defp xor(key, right_e) do
    for {x, y} <- (Enum.zip key, right_e), do: x ^^^ y 
  end

  @spec encrypt_block([integer()], integer()) :: [integer()]
  defp encrypt_block(block, keys \\ '', n \\ 0)

  defp encrypt_block(block, keys, n) when n == 0 do
    keys = generate_keys('12345678')
    block = initial_permutation block

    {left, right} = split_block(block)
    d_e = expansion(right)
    tmp = xor Enum.at(keys, n), d_e
    #tmp = permute tmp, @p
    #tmp = substitute(tmp)
    IO.inspect(tmp)
    tmp = xor left, tmp
    left = right 
    right = tmp
    block = left ++ right 
    encrypt_block(block, keys, n + 1)
  end

  defp encrypt_block(block, keys, n) when n == @num_rounds do
    {left, right} = split_block(block)
    final_permutation (right ++ left) 
  end

  defp encrypt_block(block, keys, n) do
    {left, right} = split_block(block)
    d_e = expansion(right)
    
    tmp = xor Enum.at(keys, n), d_e
    #tmp = permute tmp, @p
    #tmp = substitute(tmp)
    IO.inspect(tmp)
    tmp = xor left, tmp
    left = right 
    right = tmp
    block = left ++ right
    encrypt_block(block, keys, n + 1)
  end

  # Encrypt the given plain text, which must be a an list of blocks
  @spec encrypt_blocks([integer()]) :: [[integer()]]
  defp encrypt_blocks(plain) do
    Enum.map(plain, &encrypt_block/1)
  end

  @spec decrypt_block([integer()], [integer()], integer()) :: [integer()]
  defp decrypt_block(block, keys \\ '', n \\ 0)

  defp decrypt_block(block, keys, n) when n == 0 do
    keys = generate_keys('12345678')
    block =  initial_permutation block
    {left, right} = split_block(block)
    d_e = expansion(right)
    
    tmp = xor Enum.at(keys, 15-n), d_e
    #tmp = permute tmp, @p
    tmp = xor left, tmp
    left = right 
    right = tmp
    block = left ++ right
    decrypt_block(block, keys, n + 1)
  end

  defp decrypt_block(block, keys, n) when n == @num_rounds do
    {left, right} = split_block(block)
    final_permutation (right ++ left) 
  end

  defp decrypt_block(block, keys, n) do
    IO.puts "================="
    IO.puts n
    IO.puts "================="
    {left, right} = split_block(block)
    d_e = expansion(right)
    
    tmp = xor Enum.at(keys, 15-n), d_e
    
    #tmp = permute tmp, @p
    tmp = xor left, tmp
    left = right 
    right = tmp
    block = left ++ right
    decrypt_block(block, keys, n + 1)
  end

  @spec split_blocks([integer()]) :: [[integer()]]
  defp split_blocks(text) do
    leftover = List.duplicate (hd ' '), @block_size_bytes
    Enum.chunk_every text, @block_size_bytes, @block_size_bytes, leftover
  end

  # Descrypt the given plain text, which must be a list of blocks
  defp decrypt_blocks(cyphered) do
    Enum.map(cyphered, &decrypt_block/1)
  end

  @spec encrypt([integer()]) :: String.t()
  def encrypt(plain) do
    IO.puts "Encrypting..."
    split_blocks(plain) |> encrypt_blocks |> Enum.join
  end

  @spec decrypt([integer()]) :: String.t()
  def decrypt(plain) do
    IO.puts "Decrypting..."
    split_blocks(plain) |> decrypt_blocks |> Enum.join
  end

  @doc """
  Used for calling the module as an stand-alone file
  """
  @spec main() :: atom()
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
    :ok
  end
end

DES.main()
