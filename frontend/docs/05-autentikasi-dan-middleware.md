# Modul 05: Autentikasi NextAuth.js & Middleware Security

## 1. Arsitektur Keamanan Dashboard

Dashboard admin dan staf resepsionis hotel (`/admin/*`) dilindungi oleh sistem autentikasi terpadu menggunakan **NextAuth.js (Auth.js)** dengan strategi **JSON Web Token (JWT)**.

```
       Browser Staf                                 Next.js Server                       Go Backend
            │                                             │                                  │
     Input Kredensial                                     │                                  │
     (username, password)                                 │                                  │
            │                                             │                                  │
            ├─────── signIn("credentials") ──────────────►│                                  │
            │                                             │                                  │
            │                                             ├──── POST /api/admin/login ──────►│
            │                                             │     (username, password)         │
            │                                             │                                  │
            │                                             │◄─── 200 OK + JWT Token ──────────┤
            │                                             │     (user.id, role, name)        │
            │                                             │                                  │
            │                                      Terbitkan JWT                              │
            │                                      Session Cookie                             │
            │                                             │                                  │
            │◄────── Redirect ke /admin ──────────────────┤                                  │
            │                                             │                                  │
     Request /admin/*                                     │                                  │
            │                                             │                                  │
            ├─────── Request dengan Cookie ──────────────►│ (middleware.ts mengecek JWT)     │
            │                                             │ - Valid: Teruskan ke Dashboard   │
            │                                             │ - Tidak Valid: Redirect /login   │
```

---

## 2. Konfigurasi Provider (`src/lib/auth.ts`)

NextAuth dikonfigurasi dengan provider `CredentialsProvider`:

```typescript
export const authOptions: NextAuthOptions = {
  session: {
    strategy: "jwt",
    maxAge: 24 * 60 * 60, // Sesi berlaku 24 jam
  },
  providers: [
    CredentialsProvider({
      name: "Hotel Staff Credentials",
      credentials: {
        username: { label: "Username", type: "text" },
        password: { label: "Password", type: "password" },
      },
      async authorize(credentials) {
        const res = await fetch(`${BACKEND_URL}/api/admin/login`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            username: credentials?.username,
            password: credentials?.password,
          }),
        });

        if (!res.ok) return null;
        const data = await res.json();
        return {
          id: String(data.user.id),
          name: data.user.name,
          email: data.user.email,
          role: data.user.role,
          token: data.token,
        };
      },
    }),
  ],
  callbacks: {
    async jwt({ token, user }) {
      if (user) {
        token.role = (user as any).role;
        token.accessToken = (user as any).token;
      }
      return token;
    },
    async session({ session, token }) {
      if (session.user) {
        (session.user as any).role = token.role;
        (session.user as any).accessToken = token.accessToken;
      }
      return session;
    },
  },
  pages: {
    signIn: "/login",
  },
};
```

---

## 3. Proteksi Route Middleware (`src/middleware.ts`)

Next.js Middleware berjalan pada lapisan *Edge/Server Runtime* sebelum halaman apapun di-render:

```typescript
export { default } from "next-auth/middleware";

export const config = {
  matcher: ["/admin/:path*"],
};
```

### Mekanisme Pengamanan:
* Jika pengguna belum login mencoba membuka `/admin`, `/admin/bookings`, `/admin/takeover`, atau `/admin/rooms`, middleware secara otomatis mencegat request dan mengarahkannya (*redirect*) ke `/login?callbackUrl=...`.
* Setelah login berhasil, pengguna dikembalikan ke halaman tujuan semula.
* Akun staf bawaan yang tersedia saat instalasi:
  - **Administrator / GM**: `admin` / `admin123`
  - **Front Desk / Resepsionis**: `receptionist` / `staff123`
