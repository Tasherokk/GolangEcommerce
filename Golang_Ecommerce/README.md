# PostgreSQL Export

Generated on 2024-11-20 18:57:04.724227302 UTC by Fabricate v1.1.0

## Instructions

To load the data into a PostgreSQL database, execute the following command replacing the values with your own:

```bash
psql postgres://<user>:<password>@<host>:<port>/<db> -f load.sql -v current_dir="$(pwd)"
```

## Exported tables

This is the list of exported tables, with their corresponding row count and file names:

    public.AuditLogs: 100 rows => AuditLogs-e0af497e7235c67b817d0b1164f07db5dc62286d3aae246a6e58f7204bd0df8c.csv
    public.CartItems: 100 rows => CartItems-52e1c36deb1112c24fd882733338fa2133097017d7fa4df794dfa93ebb7e38f2.csv
    public.Categories: 10 rows => Categories-7098fc3c09edb7febea7cfb3fa01cf6392210a9b537abff9e1338f8ae2c944d8.csv
    public.OrderItems: 100 rows => OrderItems-91307b3f5ce3d8c2ba86639371ecad6506c8b4f4a205a1c776ff2438ee02f18d.csv
    public.Orders: 100 rows => Orders-77542fe84a475bbb3a44c477a588acc6f771d2cb5785fba17debbdb25249ccc9.csv
    public.Payments: 100 rows => Payments-7b63c58c6e1f13875a7961331c5fcc1ed1af45c9de8416a7d8380bbcb665a6a5.csv
    public.ProductImages: 100 rows => ProductImages-2575fedd43162dad9c224544ace11bd560e7782ab20548d15c8aec0f4ec7c566.csv
    public.Products: 100 rows => Products-31754911f0a23c1e1ad25e110a85cf04771dfe49a012347d922dc6472906999e.csv
    public.Reviews: 100 rows => Reviews-9de2bb82e71ead404287635d46bcbeb41cba67dff9110fb047abce8d6b2f284f.csv
    public.Roles: 5 rows => Roles-3015edee696e5be3597307f4f1ef06550e508d58503e50287f221e2b2c3095a6.csv
    public.Sessions: 100 rows => Sessions-c4e174f143f841e687807ba76a819d24aced3a9f74247ba4704309ec354f4af8.csv
    public.ShoppingCart: 100 rows => ShoppingCart-5bdadb3b0f117436246456daf50c0947d85e3036905b4e45f7104a047a25c5dc.csv
    public.UserAddresses: 100 rows => UserAddresses-ca063a00335d9641c6c9f9d64b66408c6a7720ebeb525b8e3d50fcd153c586a2.csv
    public.Users: 100 rows => Users-f444436cd877c9201a9f4780410c73d85841b506bae369e2c0d5c0a9937bb39a.csv