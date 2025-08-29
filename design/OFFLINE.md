# Design Jika Online

![topology](/docs/assets/topology.jpg)

Ketika pengecekan dengan timeout ke `/ping` tidak mendapat respon. Maka sistem akan menandai server sedang offline dan menyimpan data dari server saat masih online. Seluruh proses dan data seperti tarif dan transaksi akan dicatat dalam database lokal, selama kartu yang digunakan lolos dari proses validasi.

Saat sistem kembali online, `sync scheduler` akan melakukan `push` ke server. Tugas server adalah memproses batch data berdasarkan waktu dan tanggal transaksi, dan melakukan pengolahan lebih lanjut.
