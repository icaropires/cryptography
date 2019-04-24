defmodule DES do
  @block_size_bytes 8
  # @key_size_bytes 7

  require System
  require IO
  require String

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

  defp read_file(file_path) do
    plain = File.read!file_path |> String.trim
    IO.puts "============================== Plaintext ==========================="
    IO.puts plain
    IO.puts "===================================================================="
    plain
  end

  defp write_to(file_path, content) do
    File.write file_path, content
    IO.puts "Result saved to file '#{file_path}'..."
  end

  def encrypt_file(plain) do
    IO.puts "Encrypting..."

    plain
  end

  def decrypt_file(cyphered) do
    IO.puts "Decrypting..."

    cyphered
  end

  def main() do
    {operation, in_file, out_file} = parse_input()

    text = (read_file in_file) |> to_charlist
    blocks = Enum.chunk_every text, @block_size_bytes

    process_function = fn
      text when operation == :enc ->
        encrypt_file text
      text ->
        decrypt_file text
    end

    result = process_function.(blocks) |> Enum.join
    write_to out_file, result
  end
end

DES.main()
