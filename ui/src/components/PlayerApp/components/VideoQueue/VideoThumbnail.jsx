import React from "react"
import { Image } from 'antd';

export default function VideoThumbnail(props) {
  // const ytThumb = "https://i.ytimg.com/vi/VEDIO_ID/hq720.jpg"
  // const DMThumb = "https://s3.eu-central-1.amazonaws.com/centaur-wp/designweek/prod/content/uploads/2015/03/dailymotion-01-1002x564.jpg"
  // const VimeoThumb = "https://variety.com/wp-content/uploads/2014/01/vimeo_logo.jpg?w=912"
  const MP4Thumb = 'https://www.pngitem.com/pimgs/m/346-3466836_amazon-video-logo-01-circle-hd-png-download.png';

  const getIMG = (url) => {
    
    if (url === undefined || url === null) {
      return MP4Thumb;
    }
    return url;
  };
  return (
    <Image preview={false} width="78px"  src={getIMG(props.url)} />
  );
}
