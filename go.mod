module sumeru_custom_addons

go 1.26.2

replace sumeru => ../sumeru

replace sumeru_addons => ../sumeru_addons

require (
	sumeru v0.0.0
	sumeru_addons v0.0.0
)

require (
	github.com/lib/pq v1.12.3 // indirect
	golang.org/x/crypto v0.51.0 // indirect
)
