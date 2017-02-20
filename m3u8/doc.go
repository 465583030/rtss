// Package m3u8 is mux & demux library for Apple HLS.
// This package is a simple text formating and parsing library, so it must be simple too.
// It did not offer ways to play HLS or handle playlists over HTTP. Library features are:

//  * Support HLS specs up to version 5 of the protocol.
//  * Parsing and generation of master-playlists and media-playlists.
//  * Autodetect input streams as master or media playlists.
//  * Offer structures for keeping playlists metadata.
//  * Encryption keys support for usage with DRM systems like Verimatrix (http://verimatrix.com) etc.
//  * Support for non standard Google Widevine (http://www.widevine.com) tags.

// RFC: https://datatracker.ietf.org/doc/draft-pantos-http-live-streaming/

package m3u8
