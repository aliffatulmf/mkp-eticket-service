# Design Jika Online

![topology](/docs/assets/topology.jpg)

Sistem akan melakukan pengecekan ke endpoint `/ping` secara berkala untuk memasikan koneksi tetap terjalin.

Saat sistem dalam keadaan online, setiap `tap-in` dan `tap-out` yang lolos dari proses validasi, akan disimpan dalam cache atau database sementara sampai terkirim ke server. `sync scheduler` akan secara berkala melakukan singkronisasi dua arah.
