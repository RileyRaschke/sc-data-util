package util

type ProfileEntry struct {
	PriceInTicks   int
	Volume         uint
	BidVolume      uint
	AskVolume      uint
	NumberOfTrades uint
}

/*
   s_VolumeAtPriceV2();
   s_VolumeAtPriceV2
      ( const unsigned int Volume
      , const unsigned int BidVolume
      , const unsigned int AskVolume
      , const unsigned int NumberOfTrades
      );
*/
