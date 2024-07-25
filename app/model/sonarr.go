package model

import (
	"XArr-Rss/app/model/dbmodel"
	"time"
)

type SonarrSeries struct {
	Title           string `json:"title"`
	AlternateTitles []struct {
		Title string `json:"title"`
		//SeasonNumber      int    `json:"seasonNumber,omitempty"`
		//SceneSeasonNumber int    `json:"sceneSeasonNumber,omitempty"`
	} `json:"alternateTitles"`
	//SortTitle         string    `json:"sortTitle"`
	SeasonCount       int `json:"seasonCount"`
	TotalEpisodeCount int `json:"totalEpisodeCount"`
	EpisodeCount      int `json:"episodeCount"`
	//EpisodeFileCount  int       `json:"episodeFileCount"`
	//SizeOnDisk        int64     `json:"sizeOnDisk"`
	//Status            string    `json:"status"`
	Overview string `json:"overview,omitempty"`
	//PreviousAiring    time.Time `json:"previousAiring,omitempty"`
	//Network           string    `json:"network,omitempty"`
	//AirTime           string    `json:"airTime,omitempty"`
	//Images            []struct {
	//	CoverType string `json:"coverType"`
	//	Url       string `json:"url"`
	//	RemoteUrl string `json:"remoteUrl"`
	//} `json:"images"`
	Seasons []dbmodel.MediaSeason `json:"seasons"`
	Year    int                   `json:"year"`
	Path    string                `json:"path"`
	//ProfileId         int       `json:"profileId"`
	//LanguageProfileId int       `json:"languageProfileId"`
	//SeasonFolder      bool      `json:"seasonFolder"`
	Monitored bool `json:"monitored"`
	//UseSceneNumbering bool      `json:"useSceneNumbering"`
	//Runtime           int       `json:"runtime"`
	TvdbId   int `json:"tvdbId"`
	TvRageId int `json:"tvRageId"`
	TvMazeId int `json:"tvMazeId"`
	//FirstAired        time.Time `json:"firstAired,omitempty"`
	//LastInfoSync  time.Time `json:"lastInfoSync"`
	SeriesType string `json:"seriesType"`
	//CleanTitle    string    `json:"cleanTitle"`
	ImdbId    string `json:"imdbId,omitempty"`
	TitleSlug string `json:"titleSlug"`
	//Certification string    `json:"certification,omitempty"`
	//Genres        []string  `json:"genres"`
	Tags []int `json:"tags"`
	//Added         time.Time `json:"added"`
	//Ratings       struct {
	//	Votes int     `json:"votes"`
	//	Value float64 `json:"value"`
	//} `json:"ratings"`
	//QualityProfileId int       `json:"qualityProfileId"`
	Id int `json:"id"`
	//NextAiring       time.Time `json:"nextAiring,omitempty"`

	// v4
	Statistics struct {
		SeasonCount       int `json:"seasonCount"`
		EpisodeFileCount  int `json:"episodeFileCount"`
		EpisodeCount      int `json:"episodeCount"`
		TotalEpisodeCount int `json:"totalEpisodeCount"`
		//SizeOnDisk        int           `json:"sizeOnDisk"`
		//ReleaseGroups     []interface{} `json:"releaseGroups"`
		//PercentOfEpisodes int           `json:"percentOfEpisodes"`
	} `json:"statistics,omitempty"`
}

type SonarrSeriesV4 struct {
	Title           string `json:"title"`
	AlternateTitles []struct {
		Title string `json:"title"`
		//SceneSeasonNumber int    `json:"sceneSeasonNumber"`
	} `json:"alternateTitles"`
	//SortTitle      string    `json:"sortTitle"`
	//Status   string `json:"status"`
	//Ended    bool   `json:"ended"`
	Overview string `json:"overview"`
	//PreviousAiring time.Time `json:"previousAiring"`
	//Network        string    `json:"network"`
	//AirTime        string    `json:"airTime"`
	//Images         []struct {
	//	CoverType string `json:"coverType"`
	//	Url       string `json:"url"`
	//	RemoteUrl string `json:"remoteUrl"`
	//} `json:"images"`
	//OriginalLanguage struct {
	//	Id   int    `json:"id"`
	//	Name string `json:"name"`
	//} `json:"originalLanguage"`
	Seasons []dbmodel.MediaSeason `json:"seasons"`
	Year    int                   `json:"year"`
	Path    string                `json:"path"`
	//QualityProfileId  int       `json:"qualityProfileId"`
	//SeasonFolder      bool      `json:"seasonFolder"`
	//Monitored         bool      `json:"monitored"`
	//UseSceneNumbering bool      `json:"useSceneNumbering"`
	//Runtime           int       `json:"runtime"`
	TvdbId   int `json:"tvdbId"`
	TvRageId int `json:"tvRageId"`
	TvMazeId int `json:"tvMazeId"`
	//FirstAired        time.Time `json:"firstAired"`
	SeriesType string `json:"seriesType"`
	//CleanTitle        string    `json:"cleanTitle"`
	ImdbId    string `json:"imdbId"`
	TitleSlug string `json:"titleSlug"`
	//RootFolderPath    string    `json:"rootFolderPath"`
	//Certification     string    `json:"certification"`
	//Genres []string `json:"genres"`
	Tags []int `json:"tags"`
	//Added             time.Time `json:"added"`
	//Ratings           struct {
	//	Votes int `json:"votes"`
	//	Value int `json:"value"`
	//} `json:"ratings"`
	Statistics struct {
		SeasonCount       int `json:"seasonCount"`
		EpisodeFileCount  int `json:"episodeFileCount"`
		EpisodeCount      int `json:"episodeCount"`
		TotalEpisodeCount int `json:"totalEpisodeCount"`
		//SizeOnDisk        int           `json:"sizeOnDisk"`
		//ReleaseGroups     []interface{} `json:"releaseGroups"`
		//PercentOfEpisodes int           `json:"percentOfEpisodes"`
	} `json:"statistics"`
	//LanguageProfileId int       `json:"languageProfileId"`
	Id int `json:"id"`
	//NextAiring        time.Time `json:"nextAiring,omitempty"`
}

type SonarrTags []SonarrTag
type SonarrTag struct {
	Label string `json:"label"`
	Id    int    `json:"id"`
}

type SonarrSystemStatusRes struct {
	AppName string `json:"appName"`
	Version string `json:"version"`
	//BuildTime              string `json:"buildTime"`
	//IsDebug                bool   `json:"isDebug"`
	//IsProduction           bool   `json:"isProduction"`
	//IsAdmin                bool   `json:"isAdmin"`
	//IsUserInteractive      bool   `json:"isUserInteractive"`
	//StartupPath            string `json:"startupPath"`
	//AppData                string `json:"appData"`
	//OsName                 string `json:"osName"`
	//OsVersion              string `json:"osVersion"`
	//IsMonoRuntime          bool   `json:"isMonoRuntime"`
	//IsMono                 bool   `json:"isMono"`
	//IsLinux                bool   `json:"isLinux"`
	//IsOsx                  bool   `json:"isOsx"`
	//IsWindows              bool   `json:"isWindows"`
	//Mode                   string `json:"mode"`
	//Branch                 string `json:"branch"`
	//Authentication         string `json:"authentication"`
	//SqliteVersion          string `json:"sqliteVersion"`
	//UrlBase                string `json:"urlBase"`
	//RuntimeVersion         string `json:"runtimeVersion"`
	//RuntimeName            string `json:"runtimeName"`
	//StartTime              string `json:"startTime"`
	//PackageVersion         string `json:"packageVersion"`
	//PackageAuthor          string `json:"packageAuthor"`
	//PackageUpdateMechanism string `json:"packageUpdateMechanism"`
}

type SonarrEpisode struct {
	SeriesId                 int       `json:"seriesId"`
	EpisodeFileId            int       `json:"episodeFileId"`
	SeasonNumber             int       `json:"seasonNumber"`
	EpisodeNumber            int       `json:"episodeNumber"`
	Title                    string    `json:"title"`
	AirDate                  string    `json:"airDate"`
	AirDateUtc               time.Time `json:"airDateUtc"`
	HasFile                  bool      `json:"hasFile"`
	Monitored                bool      `json:"monitored"`
	AbsoluteEpisodeNumber    int       `json:"absoluteEpisodeNumber"`
	UnverifiedSceneNumbering bool      `json:"unverifiedSceneNumbering"`
	Id                       int       `json:"id"`
}

// Sonarr 下载器列表
type SonarrDownloadList struct {
	Enable                   bool   `json:"enable"`
	Protocol                 string `json:"protocol"`
	Priority                 int    `json:"priority"`
	RemoveCompletedDownloads bool   `json:"removeCompletedDownloads"`
	RemoveFailedDownloads    bool   `json:"removeFailedDownloads"`
	Name                     string `json:"name"`
	Fields                   []struct {
		Order         int         `json:"order"`
		Name          string      `json:"name"`
		Label         string      `json:"label"`
		Value         interface{} `json:"value,omitempty"`
		Type          string      `json:"type"`
		Advanced      bool        `json:"advanced"`
		HelpText      string      `json:"helpText,omitempty"`
		SelectOptions []struct {
			Value int    `json:"value"`
			Name  string `json:"name"`
			Order int    `json:"order"`
		} `json:"selectOptions,omitempty"`
	} `json:"fields"`
	ImplementationName string        `json:"implementationName"`
	Implementation     string        `json:"implementation"`
	ConfigContract     string        `json:"configContract"`
	InfoLink           string        `json:"infoLink"`
	Tags               []interface{} `json:"tags"`
	Id                 int           `json:"id"`
}
