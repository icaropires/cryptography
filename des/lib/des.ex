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
  def encrypt_blocks(plain) do
    IO.puts "Encrypting..."

    plain
  end

  @doc """
  Descrypt the given plain text, which must be a an list of blocks
  """
  def decrypt_blocks(cyphered) do
    IO.puts "Decrypting..."

    cyphered
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
