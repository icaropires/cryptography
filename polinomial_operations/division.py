def invert(a, b):
    def degree(a):
      res = 0
      a >>= 1

      while a:
        a >>= 1;
        res += 1;

      return res

    x1, x2, j = 1, 0, degree(a) - 8
    while a != 1:
        if j < 0:
          a, b = b, a
          x1, x2 = x2, x1
          j = -j

        a ^= b << j
        x1 ^= x2 << j

        j = degree(a) - degree(b)

        return x1

a, b = list(map(int, input('Insira A e B separados por espaço: ').split()))
print(invert(a, b))
