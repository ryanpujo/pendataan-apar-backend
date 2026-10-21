package models

import (
	"database/sql"
	"time"
)

// Status Tabung
type KondisiTabung string

const (
	TabungBaik        KondisiTabung = "baik"
	TabungPenyokMinor KondisiTabung = "penyok_minor" // masih aman pakai
	TabungSegelLepas  KondisiTabung = "segel_lepas"  // pin/segel terbuka, tapi tabung utuh
	TabungBocor       KondisiTabung = "bocor"        // bahaya
	TabungKaratParah  KondisiTabung = "karat_parah"  // bahaya
)

// Status Selang (Hose & Nozzle)
type KondisiSelang string

const (
	SelangBaik     KondisiSelang = "baik"
	SelangRetak    KondisiSelang = "retak_rambut" // getas/retak luar tapi belum tembus
	SelangTerlepas KondisiSelang = "terlepas"     // copot dari drat/clamp
	SelangBolong   KondisiSelang = "bolong"       // bocor saat disemprot
	SelangTersumbat KondisiSelang = "tersumbat"   // ada kotoran/hewan
)

// Status Tekanan (Manometer / Barometer)
type KondisiTekanan string

const (
	TekananHijau       KondisiTekanan = "hijau"        // normal (ideal)
	TekananMerahKiri   KondisiTekanan = "merah_kiri"   // undercharge / bocor gas dorong
	TekananMerahKanan  KondisiTekanan = "merah_kanan"  // overcharge / berisiko ledak
	TekananKacaPecah   KondisiTekanan = "kaca_pecah"   // jarum normal tapi kaca manometer retak
)

// Status Media Isi (Powder, CO2, Foam)
type KondisiIsi string

const (
	IsiBaik       KondisiIsi = "baik"
	IsiKadaluarsa KondisiIsi = "kadaluarsa" // lewat masa rekomendasi servis tapi belum beku
	IsiBeku       KondisiIsi = "beku"       // menggumpal parah (powder tidak mengalir)
	IsiKosong     KondisiIsi = "kosong"     // habis dipakai
)

// Tingkat kelayakan unit secara keseluruhan
type StatusKelayakan string

const (
	StatusSiapPakai StatusKelayakan = "siap_pakai"
	StatusPerhatian StatusKelayakan = "perhatian" // cacat minor, masih bisa semprot darurat
	StatusAfkir     StatusKelayakan = "afkir"     // rusak total / membahayakan
)

// Apar model database/sql
type Apar struct {
	ID        int64          `json:"id"`
	WargaID   int64          `json:"warga_id"`
	JenisApar string         `json:"jenis_apar"`
	Tabung    KondisiTabung  `json:"tabung"`
	Selang    KondisiSelang  `json:"selang"`
	Tekanan   KondisiTekanan `json:"tekanan"`
	Isi       KondisiIsi     `json:"isi"`
	Catatan   string         `json:"catatan"`
	Expired   sql.NullTime   `json:"expired"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// Kelayakan menghitung apakah APAR masih bisa dipakai darurat atau wajib diganti
func (a Apar) Kelayakan() StatusKelayakan {
	// Kerusakan fatal (tidak bisa atau bahaya disemprot)
	if a.Tabung == TabungBocor || a.Tabung == TabungKaratParah ||
		a.Tekanan == TekananMerahKiri || a.Tekanan == TekananMerahKanan ||
		a.Isi == IsiBeku || a.Isi == IsiKosong ||
		a.Selang == SelangBolong || a.Selang == SelangTersumbat {
		return StatusAfkir
	}

	// Masih bisa semprot tapi perlu tindakan teknis
	if a.Tabung == TabungPenyokMinor || a.Tabung == TabungSegelLepas ||
		a.Selang == SelangRetak || a.Selang == SelangTerlepas ||
		a.Tekanan == TekananKacaPecah || a.Isi == IsiKadaluarsa {
		return StatusPerhatian
	}

	return StatusSiapPakai
}

type CreateAparRequest struct {
	WargaID   int64      `json:"warga_id" validate:"required"`
	JenisApar string     `json:"jenis_apar" validate:"required"`
	Tabung    string     `json:"tabung" validate:"required,oneof=baik penyok_minor segel_lepas bocor karat_parah"`
	Selang    string     `json:"selang" validate:"required,oneof=baik retak_rambut terlepas bolong tersumbat"`
	Tekanan   string     `json:"tekanan" validate:"required,oneof=hijau merah_kiri merah_kanan kaca_pecah"`
	Isi       string     `json:"isi" validate:"required,oneof=baik kadaluarsa beku kosong"`
	Catatan   string     `json:"catatan,omitempty"`
	Expired   *time.Time `json:"expired,omitempty"`
}

type UpdateAparRequest struct {
	WargaID   *int64     `json:"warga_id,omitempty"`
	JenisApar *string    `json:"jenis_apar,omitempty"`
	Tabung    *string    `json:"tabung,omitempty" validate:"omitempty,oneof=baik penyok_minor segel_lepas bocor karat_parah"`
	Selang    *string    `json:"selang,omitempty" validate:"omitempty,oneof=baik retak_rambut terlepas bolong tersumbat"`
	Tekanan   *string    `json:"tekanan,omitempty" validate:"omitempty,oneof=hijau merah_kiri merah_kanan kaca_pecah"`
	Isi       *string    `json:"isi,omitempty" validate:"omitempty,oneof=baik kadaluarsa beku kosong"`
	Catatan   *string    `json:"catatan,omitempty"`
	Expired   *time.Time `json:"expired,omitempty"`
}