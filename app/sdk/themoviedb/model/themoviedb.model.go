package model

type ApiSearchMultiRes struct {
	Page         int                    `json:"page"`
	Results      []ApiSearchMultiResult `json:"results"`
	TotalPages   int                    `json:"total_pages"`
	TotalResults int                    `json:"total_results"`
}

type ApiSearchMultiResult struct {
	Adult            bool    `json:"adult"`
	BackdropPath     string  `json:"backdrop_path"`
	GenreIds         []int   `json:"genre_ids"`
	Id               int     `json:"id"`
	MediaType        string  `json:"media_type"`
	OriginalLanguage string  `json:"original_language"`
	Overview         string  `json:"overview"`
	Popularity       float64 `json:"popularity"`
	PosterPath       string  `json:"poster_path"`
	ReleaseDate      string  `json:"release_date"`
	Video            bool    `json:"video"`
	VoteAverage      float64 `json:"vote_average"`
	VoteCount        int     `json:"vote_count"`
	// tv 才有下面的东西
	FirstAirDate  string   `json:"first_air_date"`
	Name          string   `json:"name"`
	OriginCountry []string `json:"origin_country"`
	OriginalName  string   `json:"original_name"`

	// 电影才有的
	Title         string `json:"title"`
	OriginalTitle string `json:"original_title"`
	ImdbId        string `json:"imdb_id"`
}

type ApiMovieInfoRes struct {
	Adult               bool        `json:"adult"`
	BackdropPath        string      `json:"backdrop_path"`
	BelongsToCollection interface{} `json:"belongs_to_collection"`
	Budget              int         `json:"budget"`
	Genres              []struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
	Homepage            string  `json:"homepage"`
	Id                  int     `json:"id"`
	ImdbId              string  `json:"imdb_id"`
	OriginalLanguage    string  `json:"original_language"`
	OriginalTitle       string  `json:"original_title"`
	Overview            string  `json:"overview"`
	Popularity          float64 `json:"popularity"`
	PosterPath          string  `json:"poster_path"`
	ProductionCompanies []struct {
		Id            int     `json:"id"`
		LogoPath      *string `json:"logo_path"`
		Name          string  `json:"name"`
		OriginCountry string  `json:"origin_country"`
	} `json:"production_companies,omitempty"`
	ProductionCountries []struct {
		Iso31661 string `json:"iso_3166_1"`
		Name     string `json:"name"`
	} `json:"production_countries,omitempty"`
	ReleaseDate     string `json:"release_date"`
	Revenue         int    `json:"revenue"`
	Runtime         int    `json:"runtime"`
	SpokenLanguages []struct {
		EnglishName string `json:"english_name"`
		Iso6391     string `json:"iso_639_1"`
		Name        string `json:"name"`
	} `json:"spoken_languages"`
	Status      string  `json:"status"`
	Tagline     string  `json:"tagline"`
	Title       string  `json:"title"`
	Video       bool    `json:"video"`
	VoteAverage float64 `json:"vote_average"`
	VoteCount   int     `json:"vote_count"`
}

type ApiGenresRes struct {
	Genres []struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
}

type ApiTvExtendIds struct {
	Id          int         `json:"id"`
	ImdbId      string      `json:"imdb_id,omitempty"`
	FreebaseMid interface{} `json:"freebase_mid,omitempty"`
	FreebaseId  interface{} `json:"freebase_id,omitempty"`
	TvdbId      int         `json:"tvdb_id,omitempty"`
	TvrageId    interface{} `json:"tvrage_id,omitempty"`
	FacebookId  string      `json:"facebook_id,omitempty"`
	InstagramId string      `json:"instagram_id,omitempty"`
	TwitterId   string      `json:"twitter_id,omitempty"`
}

type ApiFindByExtendIds struct {
	MovieResults     []ApiFindByExtendIdsMovieResult `json:"movie_results"`
	PersonResults    []interface{}                   `json:"person_results"`
	TvResults        []ApiFindByExtendIdsTvResult    `json:"tv_results"`
	TvEpisodeResults []interface{}                   `json:"tv_episode_results"`
	TvSeasonResults  []interface{}                   `json:"tv_season_results"`
}

type ApiFindByExtendIdsMovieResult struct {
	Adult            bool    `json:"adult"`
	BackdropPath     string  `json:"backdrop_path"`
	GenreIds         []int   `json:"genre_ids"`
	Id               int     `json:"id"`
	OriginalLanguage string  `json:"original_language"`
	OriginalTitle    string  `json:"original_title"`
	Overview         string  `json:"overview"`
	PosterPath       string  `json:"poster_path"`
	ReleaseDate      string  `json:"release_date"`
	Title            string  `json:"title"`
	Video            bool    `json:"video"`
	VoteAverage      float64 `json:"vote_average"`
	VoteCount        int     `json:"vote_count"`
	Popularity       float64 `json:"popularity"`
}

type ApiFindByExtendIdsTvResult struct {
	BackdropPath     string   `json:"backdrop_path"`
	FirstAirDate     string   `json:"first_air_date"`
	GenreIds         []int    `json:"genre_ids"`
	Id               int      `json:"id"`
	Name             string   `json:"name"`
	OriginCountry    []string `json:"origin_country"`
	OriginalLanguage string   `json:"original_language"`
	OriginalName     string   `json:"original_name"`
	Overview         string   `json:"overview"`
	PosterPath       string   `json:"poster_path"`
	VoteAverage      float64  `json:"vote_average"`
	VoteCount        int      `json:"vote_count"`
	Popularity       float64  `json:"popularity"`
}

type ApiSearchMovieResult struct {
	Page    int `json:"page"`
	Results []struct {
		Adult            bool    `json:"adult"`
		BackdropPath     string  `json:"backdrop_path"`
		GenreIds         []int   `json:"genre_ids"`
		Id               int     `json:"id"`
		OriginalLanguage string  `json:"original_language"`
		OriginalTitle    string  `json:"original_title"`
		Overview         string  `json:"overview"`
		Popularity       float64 `json:"popularity"`
		PosterPath       string  `json:"poster_path"`
		ReleaseDate      string  `json:"release_date"`
		Title            string  `json:"title"`
		Video            bool    `json:"video"`
		VoteAverage      float64 `json:"vote_average"`
		VoteCount        int     `json:"vote_count"`
	} `json:"results"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

type ApiSearchTvResult struct {
	Page    int `json:"page"`
	Results []struct {
		BackdropPath     string   `json:"backdrop_path"`
		FirstAirDate     string   `json:"first_air_date"`
		GenreIds         []int    `json:"genre_ids"`
		Id               int      `json:"id"`
		Name             string   `json:"name"`
		OriginCountry    []string `json:"origin_country"`
		OriginalLanguage string   `json:"original_language"`
		OriginalName     string   `json:"original_name"`
		Overview         string   `json:"overview"`
		Popularity       float64  `json:"popularity"`
		PosterPath       string   `json:"poster_path"`
		VoteAverage      float64  `json:"vote_average"`
		VoteCount        int      `json:"vote_count"`
	} `json:"results"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

type ApiMovieTopRate ApiMovieUpcoming

type ApiMovieUpcoming struct {
	Dates struct {
		Maximum string `json:"maximum"`
		Minimum string `json:"minimum"`
	} `json:"dates"`
	Page    int `json:"page"`
	Results []struct {
		Adult            bool    `json:"adult"`
		BackdropPath     string  `json:"backdrop_path"`
		GenreIds         []int   `json:"genre_ids"`
		Id               int     `json:"id"`
		OriginalLanguage string  `json:"original_language"`
		OriginalTitle    string  `json:"original_title"`
		Overview         string  `json:"overview"`
		Popularity       float64 `json:"popularity"`
		PosterPath       string  `json:"poster_path"`
		ReleaseDate      string  `json:"release_date"`
		Title            string  `json:"title"`
		Video            bool    `json:"video"`
		VoteAverage      float64 `json:"vote_average"`
		VoteCount        int     `json:"vote_count"`
	} `json:"results"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

type ApiTvTopRate struct {
	Page    int `json:"page"`
	Results []struct {
		BackdropPath     string   `json:"backdrop_path"`
		FirstAirDate     string   `json:"first_air_date"`
		GenreIds         []int    `json:"genre_ids"`
		Id               int      `json:"id"`
		Name             string   `json:"name"`
		OriginCountry    []string `json:"origin_country"`
		OriginalLanguage string   `json:"original_language"`
		OriginalName     string   `json:"original_name"`
		Overview         string   `json:"overview"`
		Popularity       float64  `json:"popularity"`
		PosterPath       string   `json:"poster_path"`
		VoteAverage      float64  `json:"vote_average"`
		VoteCount        int      `json:"vote_count"`
	} `json:"results"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

type ApiTvPopular ApiTvTopRate

type ApiConfigurationRs struct {
	Images struct {
		BaseUrl       string   `json:"base_url"`
		SecureBaseUrl string   `json:"secure_base_url"`
		BackdropSizes []string `json:"backdrop_sizes"`
		LogoSizes     []string `json:"logo_sizes"`
		PosterSizes   []string `json:"poster_sizes"`
		ProfileSizes  []string `json:"profile_sizes"`
		StillSizes    []string `json:"still_sizes"`
	} `json:"images"`
	ChangeKeys []string `json:"change_keys"`
}

type ApiAlternativeTitles struct {
	Id      int                          `json:"id"`
	Results []ApiAlternativeTitlesResult `json:"results"`
}

type ApiAlternativeTitlesResult struct {
	Iso31661 string `json:"iso_3166_1"`
	Title    string `json:"title"`
	Type     string `json:"type"`
}
