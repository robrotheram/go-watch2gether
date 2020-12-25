import React from "react"
import ReactPlayer from 'react-player'
import {connect} from 'react-redux'
import {play, pause,updateSeek, updateQueue, handleFinish} from '../store/room/room.actions'
import { Empty } from "antd"

class VideoPlayer extends React.Component {
    state = {
        url: null,
        pip: false,
        playing: true,
        light: false,
        volume: 0.8,
        muted: false,
        played: 0,
        loaded: 0,
        duration: 0,
        playbackRate: 1.0,
        loop: false
    }


    componentWillReceiveProps(nextProps) {
        // You don't have to do this check first, but it can help prevent an unneeded render
        if (nextProps.seek !== this.state.played) {
            console.log("STATE_CHANGE",nextProps )
          this.setState({ played: nextProps.seek });
          if (this.player !== undefined) {
            this.player.seekTo(parseFloat(nextProps.seek))
          }
        }
      }

    componentDidMount() {
        let _this = this
        setTimeout(function(){
            if (_this.player !== undefined && _this.player !== null){
                _this.player.seekTo(parseFloat(_this.props.seek)) 
            }
        }, 200);
        
    }

    handlePlay = () => {
        if(!this.props.playing){
            play();
        }
        
    }

    handlePause = () => {
        if(this.props.playing){
            pause();
        }
    }

    handleEnded = () => {
        updateSeek(1)
        handleFinish()
    }

    handleProgress = state => {
        // console.log('onProgress', state)
        // console.log("PLAYER!", this.player)
        this.setState({ seek: state.playedSeconds })
        updateSeek(state.played)
        
    }

    onProgressUpdate = e => {
        this.setState({ seeking: false })
        this.player.seekTo(parseFloat(e.target.value))
    }

    
    ref = player => { this.player = player }

    render(){
            const {playing, url } = this.props
            return(
            <div style={{ "height":"600px", "width":"100%"}}>  
            {url !== "" ? 
               <ReactPlayer 
                ref={this.ref}
                width="100%" height="600px"  
                url={url} 
                controls={true}
                playing={playing}
                onPause={this.handlePause}
                onPlay={this.handlePlay}
                onProgress={this.handleProgress}
                onEnded={this.handleEnded}
                />  
            : <Empty  style={{ "height":"600px", "width":"100%", "paddingTop":"180px"}}/>}
          </div>  
        )
    }
}
const mapStateToProps  = (state) =>{
    return state.video
  } 
export default connect(mapStateToProps, {updateQueue, handleFinish, updateSeek})(VideoPlayer)
  