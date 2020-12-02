import { Image } from 'antd';
export function VideoThumbnail(props) {
    // const ytThumb = "https://i.ytimg.com/vi/VEDIO_ID/hq720.jpg"
    // const DMThumb = "https://s3.eu-central-1.amazonaws.com/centaur-wp/designweek/prod/content/uploads/2015/03/dailymotion-01-1002x564.jpg"
    // const VimeoThumb = "https://variety.com/wp-content/uploads/2014/01/vimeo_logo.jpg?w=912"
    const MP4Thumb = "https://www.pngitem.com/pimgs/m/346-3466836_amazon-video-logo-01-circle-hd-png-download.png"
  

    function getYTThumb (url, size) {
      if (url === null) {
          return '';
      }
      size    = (size === null) ? 'big' : size;
      let results = url.match('[\\?&]v=([^&#]*)');
      let video   = (results === null) ? url : results[1];

      if (size === 'small') {
          return 'http://img.youtube.com/vi/' + video + '/2.jpg';
      }
      return 'http://img.youtube.com/vi/' + video + '/mqdefault.jpg';
    };

    function youtube_parser(url){
      var regExp = /^.*((youtu.be\/)|(v\/)|(\/u\/\w\/)|(embed\/)|(watch\?))\??v?=?([^#]*).*/;
      var match = url.match(regExp);
      if (match&&match[7].length===11){
          return true
      }else{
          return false
      }
    }


    const getIMG = (url) => {
      if( url === undefined || url === null ){
        return MP4Thumb
      }
      if (youtube_parser(url)) {
        return getYTThumb(url, 'big')
      }
      return MP4Thumb
    }
    return(
        <Image width={100} height={60} src={getIMG(props.url)} />
      )
  }