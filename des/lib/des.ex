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

  @doc """
  Encrypt the given plain text, which must be a an list of blocks
  """
  def encrypt_blocks(plain, n\\0) do
    if n > 5 do
      IO.puts("42424242")
      IO.inspect(plain)
      IO.puts("42424242")
      {:ok, file} = File.open("newfile")
      File.write("aisim", to_string(plain))
      plain 
    else
      IO.puts "Encrypting..."

      IO.inspect(plain)
      list = (List.to_string(to_charlist(plain)))

      left = String.slice(list, 0..(div(String.length(list),2)-1))
      right = String.slice(list, (div(String.length(list),2))..-1)

      right1 = Enum.map(to_charlist(right), fn x -> x + 3 end )
      c = for {x, y} <- (Enum.zip to_charlist(left), to_charlist(right1)), do: x ^^^ y 
      IO.inspect(c)

      list = to_charlist(right) ++ c
      IO.puts("--------------------")
      IO.inspect(list)
      IO.puts("--------------------")
      encrypt_blocks(list , n + 1)
    end
  end

  @doc """
  Descrypt the given plain text, which must be a an list of blocks
  """
  def decrypt_blocks(cyphered, n \\ 0) do
    if n > 6 do
      list = (List.to_string(to_charlist(cyphered)))
      left = String.slice(list, 0..(div(String.length(list),2)-1))
      right = String.slice(list, (div(String.length(list),2))..-1)
      list = to_charlist(right) ++ to_charlist(left)


      File.write("aisimd", to_string(list))
      cyphered 
    else
      IO.puts "Decrypting..."
      if n == 0 do
        list = (List.to_string(to_charlist(cyphered)))
        IO.inspect(list)
        left = String.slice(list, 0..(div(String.length(list),2)-1))
        right = String.slice(list, (div(String.length(list),2))..-1)
        list = to_charlist(right) ++ to_charlist(left)
        decrypt_blocks(list, n+ 1)
      else
        list = (List.to_string(to_charlist(cyphered)))
        IO.inspect(list)

        left = String.slice(list, 0..(div(String.length(list),2)-1))
        right = String.slice(list, (div(String.length(list),2))..-1)

        right1 = Enum.map(to_charlist(right), fn x -> x + 3 end )
        c = for {x, y} <- (Enum.zip to_charlist(left), to_charlist(right1)), do: x ^^^ y 
        #c = Enum.map(c, fn x -> x + 42 end )

        list = to_charlist(right)  ++ c
        IO.puts("============")
        IO.inspect(list)
        IO.puts("============")
        decrypt_blocks(list , n + 1)
      end
    end
  end

  @doc """
  Used for calling the module as an stand-alone file
  """
  def main() do
    {operation, in_file, out_file} = parse_input()

    text = (read_file in_file) |> to_charlist
    blocks = Enum.chunk_every text, @block_size_bytes

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
