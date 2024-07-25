package themoviedb

import (
	"XArr-Rss/app/sdk/themoviedb/model"
	"XArr-Rss/util/helper"
	"XArr-Rss/util/logsys"
	"encoding/json"
	"log"
	"net/url"
)

type themoviedbSdk struct {
	apiKey string
	proxy  string
}

//http://api.themoviedb.org/3/search/multi?api_key=&language=zh-CN&query=Slow%20Loop&page=1&include_adult=false

func GetTheMoviedbSdk(apikey string, proxy string) *themoviedbSdk {
	return &themoviedbSdk{apiKey: apikey, proxy: proxy}
}

func (this *themoviedbSdk) SearchMulti(query string) (error, *model.ApiSearchMultiRes) {
	uri := "http://api.themoviedb.org/3/search/multi?api_key=" + this.apiKey + "&language=zh-CN&query=" + url.QueryEscape(query) + "&page=1&include_adult=false&append_to_response=external_ids"
	//err, result := helper.CurlHelper{}.GetUri(uri, nil, []helper.CurlHeader{})
	err, result := this.api(uri, "SearchMulti."+query)
	logsys.Debug("请求影片数据接口地址:%v", "themoviedb", uri)
	if err != nil {
		log.Println("请求推送错误", err)
		return err, nil
	}
	var res model.ApiSearchMultiRes
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Println("搜索内容失败", err)
		return err, nil
	}

	log.Println("查询影片结果", string(result))
	return nil, &res
}

// 搜索电影
func (this *themoviedbSdk) SearchMovie(query string) (error, *model.ApiSearchMovieResult) {
	uri := "http://api.themoviedb.org/3/search/movie?api_key=" + this.apiKey + "&language=zh-CN&query=" + url.QueryEscape(query) + "&page=1&include_adult=false&append_to_response=external_ids"
	err, result := this.api(uri, "SearchMovie"+query)

	//err, result := helper.CurlHelper{}.GetUri(uri, nil, []helper.CurlHeader{})
	logsys.Debug("请求影片数据接口地址:%v", "themoviedb", uri)
	if err != nil {
		log.Println("请求推送错误", err)
		return err, nil
	}
	var res model.ApiSearchMovieResult
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Println("搜索内容失败", err)
		return err, nil
	}

	log.Println("查询影片结果", string(result))
	return nil, &res
}

func (this *themoviedbSdk) SearchTv(query string) (error, *model.ApiSearchTvResult) {
	uri := "http://api.themoviedb.org/3/search/tv?api_key=" + this.apiKey + "&language=zh-CN&query=" + url.QueryEscape(query) + "&page=1&include_adult=false"

	err, result := this.api(uri, "SearchTv."+query)

	//err, result := helper.CurlHelper{}.GetUri(uri, nil, []helper.CurlHeader{})
	if err != nil {
		log.Println("请求推送错误", err)
		return err, nil
	}
	var res model.ApiSearchTvResult
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Println("搜索内容失败", err)
		return err, nil
	}

	log.Println("查询影片结果", string(result))
	return nil, &res
}

func (this *themoviedbSdk) MovieById(id string) (error, *model.ApiMovieInfoRes) {
	uri := "http://api.themoviedb.org/3/movie/" + id + "?api_key=" + this.apiKey + "&language=zh-CN"
	err, result := this.api(uri, "MovieById."+id)

	//err, result := helper.CurlHelper{}.GetUri(uri, nil, []helper.CurlHeader{})
	if err != nil {
		log.Println("请求推送错误", err)
		return err, nil
	}
	var res model.ApiMovieInfoRes
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Println("搜索内容失败", err)
		return err, nil
	}

	log.Println("查询影片详情结果", string(result))
	return nil, &res
}

func (this *themoviedbSdk) GenreMovie() (error, *model.ApiGenresRes) {
	uri := "http://api.themoviedb.org/3/genre/movie/list?api_key=" + this.apiKey + "&language=zh-CN"
	err, result := this.api(uri, "GenreMovie")
	//err, result := helper.CurlHelper{}.GetUri(uri, nil, []helper.CurlHeader{})
	log.Println("请求movie标签接口地址", uri)
	if err != nil {
		log.Println("请求推送错误", err)
		return err, nil
	}
	var res model.ApiGenresRes
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Println("搜索内容失败", err)
		return err, nil
	}

	log.Println("查询movie标签详情结果", string(result))
	return nil, &res
}

func (this *themoviedbSdk) GenreTV() (error, *model.ApiGenresRes) {
	uri := "http://api.themoviedb.org/3/genre/tv/list?api_key=" + this.apiKey + "&language=zh-CN"
	log.Println("请求TV标签接口地址", uri)
	err, result := this.api(uri, "GenreTV")
	//err, result := helper.CurlHelper{}.GetUri(uri, nil, []helper.CurlHeader{})
	if err != nil {
		log.Println("请求推送错误", err)
		return err, nil
	}
	log.Println("查询TV标签详情结果", string(result))

	var res model.ApiGenresRes
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Println("搜索内容失败", err)
		return err, nil
	}

	return nil, &res
}

func (this *themoviedbSdk) TvExternalIds(tvId string) (error, *model.ApiTvExtendIds) {
	uri := "http://api.themoviedb.org/3/tv/" + tvId + "/external_ids?api_key=" + this.apiKey + "&language=zh-CN"
	err, result := this.api(uri, "TvExternalIds."+tvId)
	//err, result := helper.CurlHelper{}.GetUri(uri, nil, []helper.CurlHeader{})
	if err != nil {
		log.Println("请求推送错误", err)
		return err, nil
	}
	log.Println("请求获取额外ID详情结果", string(result))

	var res model.ApiTvExtendIds
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Println("搜索内容失败", err)
		return err, nil
	}

	return nil, &res
}
func (this *themoviedbSdk) FindByExternalIds(externalId string, externalSource string) (error, *model.ApiFindByExtendIds) {
	//Allowed Values: imdb_id, freebase_mid, freebase_id, tvdb_id, tvrage_id, facebook_id, twitter_id, instagram_id
	uri := "http://api.themoviedb.org/3/find/" + externalId + "?api_key=" + this.apiKey + "&language=zh-CN&external_source=" + externalSource
	//logsys.Debug("获取影片信息:%v", "themoviedb", uri)
	err, result := this.api(uri, "FindTvExternalIds."+externalId+externalSource)

	//err, result := helper.CurlHelper{}.GetUri(uri, nil, []helper.CurlHeader{})
	if err != nil {
		logsys.Debug("查询失败:%v", "themoviedb", err.Error())
		return err, nil
	}

	var res model.ApiFindByExtendIds
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Println("搜索内容失败", err)
		return err, nil
	}

	return nil, &res
}

func (this *themoviedbSdk) AlternativeTvTitles(tmdbId string) (error, *model.ApiAlternativeTitles) {
	//Allowed Values: imdb_id, freebase_mid, freebase_id, tvdb_id, tvrage_id, facebook_id, twitter_id, instagram_id
	uri := "http://api.themoviedb.org/3/tv/" + tmdbId + "/alternative_titles?api_key=" + this.apiKey + "&country=ZH"
	//logsys.Debug("获取影片信息:%v", "themoviedb", uri)
	err, result := this.api(uri, "AlternativeTvTitles."+tmdbId)

	//err, result := helper.CurlHelper{}.GetUri(uri, nil, []helper.CurlHeader{})
	if err != nil {
		logsys.Debug("查询失败:%v", "themoviedb", err.Error())
		return err, nil
	}

	var res model.ApiAlternativeTitles
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Println("搜索内容失败", err)
		return err, nil
	}

	return nil, &res
}

// 获取即将上映的
func (this *themoviedbSdk) GetMovieUpcoming() (err error, res *model.ApiMovieUpcoming) {
	//Allowed Values: imdb_id, freebase_mid, freebase_id, tvdb_id, tvrage_id, facebook_id, twitter_id, instagram_id
	uri := "http://api.themoviedb.org/3/movie/upcoming?api_key=" + this.apiKey + "&language=zh-CN"
	err, result := this.api(uri, "GetMovieUpcoming")
	res = &model.ApiMovieUpcoming{}
	err = json.Unmarshal(result, res)
	if err != nil {
		log.Println("获取电影即将上映失败", err)
		return err, nil
	}

	return nil, res
}

// 获取最高评价的
func (this *themoviedbSdk) GetMovieTopRate() (error, *model.ApiMovieTopRate) {
	//Allowed Values: imdb_id, freebase_mid, freebase_id, tvdb_id, tvrage_id, facebook_id, twitter_id, instagram_id
	uri := "http://api.themoviedb.org/3/movie/top_rated?api_key=" + this.apiKey + "&language=zh-CN"
	logsys.Debug("获取影片信息:%v", "themoviedb", uri)
	err, result := this.api(uri, "GetMovieTopRate")

	//err, result := helper.CurlHelper{}.GetUri(uri, nil, []helper.CurlHeader{})
	if err != nil {
		logsys.Debug("查询失败:%v", "themoviedb", err.Error())
		return err, nil
	}

	var res model.ApiMovieTopRate
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Println("获取电影最高评价的失败", err)
		return err, nil
	}

	return nil, &res
}

// 获取最喜欢的
func (this *themoviedbSdk) GetMoviePopular() (error, *model.ApiMovieTopRate) {
	//Allowed Values: imdb_id, freebase_mid, freebase_id, tvdb_id, tvrage_id, facebook_id, twitter_id, instagram_id
	uri := "http://api.themoviedb.org/3/movie/popular?api_key=" + this.apiKey + "&language=zh-CN"
	logsys.Debug("获取影片信息:%v", "themoviedb", uri)
	//err, result := helper.CurlHelper{}.GetUri(uri, nil, []helper.CurlHeader{})
	err, result := this.api(uri, "GetMoviePopular")

	if err != nil {
		logsys.Debug("查询失败:%v", "themoviedb", err.Error())
		return err, nil
	}
	//logsys.Debug("查询成功:%v", "themoviedb", string(result))

	var res model.ApiMovieTopRate
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Println("获取电影最喜欢失败", err)
		return err, nil
	}

	return nil, &res
}

// 获取正在播放的
func (this *themoviedbSdk) GetMovieNowPlaying() (error, *model.ApiMovieTopRate) {
	//Allowed Values: imdb_id, freebase_mid, freebase_id, tvdb_id, tvrage_id, facebook_id, twitter_id, instagram_id
	uri := "http://api.themoviedb.org/3/movie/now_playing?api_key=" + this.apiKey + "&language=zh-CN"
	logsys.Debug("获取影片信息:%v", "themoviedb", uri)
	err, result := this.api(uri, "GetMovieNowPlaying")
	if err != nil {
		logsys.Debug("查询失败:%v", "themoviedb", err.Error())
		return err, nil
	}
	logsys.Debug("查询成功:%v", "themoviedb", string(result))

	var res model.ApiMovieTopRate
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Println("获取电影正在播放失败", err)
		return err, nil
	}

	return nil, &res
}

func (this *themoviedbSdk) GetMovieLatest() (error, *model.ApiMovieTopRate) {
	//Allowed Values: imdb_id, freebase_mid, freebase_id, tvdb_id, tvrage_id, facebook_id, twitter_id, instagram_id
	uri := "http://api.themoviedb.org/3/movie/latest?api_key=" + this.apiKey + "&language=zh-CN"
	logsys.Debug("获取影片信息:%v", "themoviedb", uri)
	err, result := this.api(uri, "GetMovieLatest")
	if err != nil {
		logsys.Debug("查询失败:%v", "themoviedb", err.Error())
		return err, nil
	}
	logsys.Debug("查询成功:%v", "themoviedb", string(result))

	var res model.ApiMovieTopRate
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Println("获取电影最新失败", err)
		return err, nil
	}

	return nil, &res
}

func (this *themoviedbSdk) GetTvTopRate() (error, *model.ApiTvTopRate) {
	//Allowed Values: imdb_id, freebase_mid, freebase_id, tvdb_id, tvrage_id, facebook_id, twitter_id, instagram_id
	uri := "http://api.themoviedb.org/3/tv/top_rated?api_key=" + this.apiKey + "&language=zh-CN"
	logsys.Debug("获取影片信息:%v", "themoviedb", uri)
	err, result := this.api(uri, "GetTvTopRate")

	if err != nil {
		logsys.Debug("查询失败:%v", "themoviedb", err.Error())
		return err, nil
	}
	logsys.Debug("查询成功:%v", "themoviedb", string(result))

	var res model.ApiTvTopRate
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Println("获取电影最新失败", err)
		return err, nil
	}

	return nil, &res
}

func (this *themoviedbSdk) GetTvPopular() (error, *model.ApiTvPopular) {
	//Allowed Values: imdb_id, freebase_mid, freebase_id, tvdb_id, tvrage_id, facebook_id, twitter_id, instagram_id
	uri := "http://api.themoviedb.org/3/tv/popular?api_key=" + this.apiKey + "&language=zh-CN"
	logsys.Debug("获取影片信息:%v", "themoviedb", uri)
	err, result := this.api(uri, "GetTvPopular")

	if err != nil {
		logsys.Debug("查询失败:%v", "themoviedb", err.Error())
		return err, nil
	}
	logsys.Debug("查询成功:%v", "themoviedb", string(result))

	var res model.ApiTvPopular
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Println("获取电影最新失败", err)
		return err, nil
	}

	return nil, &res
}

func (this *themoviedbSdk) GetTvLatest() (error, *model.ApiTvPopular) {
	//Allowed Values: imdb_id, freebase_mid, freebase_id, tvdb_id, tvrage_id, facebook_id, twitter_id, instagram_id
	uri := "http://api.themoviedb.org/3/tv/latest?api_key=" + this.apiKey + "&language=zh-CN"
	logsys.Debug("获取影片信息:%v", "themoviedb", uri)
	err, result := this.api(uri, "GetTvLatest")

	if err != nil {
		logsys.Debug("查询失败:%v", "themoviedb", err.Error())
		return err, nil
	}
	logsys.Debug("查询成功:%v", "themoviedb", string(result))

	var res model.ApiTvPopular
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Println("获取电影最新失败", err)
		return err, nil
	}

	return nil, &res
}

func (this *themoviedbSdk) GetTvAiringToday() (error, *model.ApiTvPopular) {
	//Allowed Values: imdb_id, freebase_mid, freebase_id, tvdb_id, tvrage_id, facebook_id, twitter_id, instagram_id
	uri := "http://api.themoviedb.org/3/tv/airing_today?api_key=" + this.apiKey + "&language=zh-CN"
	logsys.Debug("获取影片信息:%v", "themoviedb", uri)
	err, result := this.api(uri, "GetTvAiringToday")

	if err != nil {
		logsys.Debug("查询失败:%v", "themoviedb", err.Error())
		return err, nil
	}

	logsys.Debug("查询成功:%v", "themoviedb", string(result))

	var res model.ApiTvPopular
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Println("获取电影最新失败", err)
		return err, nil
	}

	return nil, &res
}

func (this *themoviedbSdk) GetTvOnTheAir() (error, *model.ApiTvPopular) {
	//Allowed Values: imdb_id, freebase_mid, freebase_id, tvdb_id, tvrage_id, facebook_id, twitter_id, instagram_id
	uri := "http://api.themoviedb.org/3/tv/on_the_air?api_key=" + this.apiKey + "&language=zh-CN"
	logsys.Debug("获取影片信息:%v", "themoviedb", uri)
	err, result := this.api(uri, "GetTvOnTheAir")

	if err != nil {
		logsys.Debug("查询失败:%v", "themoviedb", err.Error())
		return err, nil
	}
	logsys.Debug("查询成功:%v", "themoviedb", string(result))

	var res model.ApiTvPopular
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Println("获取电影最新失败", err)
		return err, nil
	}

	return nil, &res
}

func (this *themoviedbSdk) api(uri string, method string) (err error, result []byte) {
	//logsys.Debug("获取影片信息(%v):%v", "themoviedb", method, uri)

	// 重新请求数据
	err, result = helper.CurlHttpHelper{}.GetProxyResult(uri, this.proxy)
	if err != nil {
		logsys.Debug("查询失败:%v", "themoviedb."+method, err.Error())
		return err, result
	}
	//
	//err, result, _ = helper.CurlHelper{}.GetUri(uri, nil, []helper.CurlHeader{}, true)
	//if err != nil {
	//	logsys.Debug("查询失败:%v", "themoviedb."+method, err.Error())
	//	return err, result
	//}
	logsys.Debug("查询成功:%v", "themoviedb."+method, string(result))
	//logsys.Debug("查询成功", "themoviedb."+method)

	// 判断是否使用缓存
	return nil, result

}

func (this *themoviedbSdk) GetConfigurationRs() (error, *model.ApiConfigurationRs) {
	//Allowed Values: imdb_id, freebase_mid, freebase_id, tvdb_id, tvrage_id, facebook_id, twitter_id, instagram_id
	uri := "http://api.themoviedb.org/3/configuration?api_key=" + this.apiKey + "&language=zh-CN"
	err, result := this.api(uri, "GetTvOnTheAir")

	if err != nil {
		logsys.Debug("查询失败:%v", "themoviedb", err.Error())
		return err, nil
	}
	var res model.ApiConfigurationRs
	err = json.Unmarshal(result, &res)
	if err != nil {
		log.Println("获取失败", err)
		return err, nil
	}
	return nil, &res
}
