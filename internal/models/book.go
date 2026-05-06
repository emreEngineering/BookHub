// Bu dosyanın models paketine ait olduğunu belirtir.
package models

// Book adında kendi veri tipimizi oluştururuz.
type Book struct {
	// ID, kitabın benzersiz numarasıdır.
	ID int `json:"id"`
	// Title, kitabın adını tutar.
	Title string `json:"title"`
	// Author, kitabın yazarını tutar.
	Author string `json:"author"`
	// Year, kitabın yayın yılını tutar.
	Year int `json:"year"`
}
