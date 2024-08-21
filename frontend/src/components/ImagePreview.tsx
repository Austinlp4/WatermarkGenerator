import React from 'react';
import { Box, CircularProgress, Typography } from '@mui/material';
import { keyframes } from '@emotion/react';

interface ImagePreviewProps {
  preview: string | null;
  isLoading: boolean;
  watermarkedImage: string | null;
}

const waveAnimation = keyframes`
  0% {
    transform: translateX(0) translateZ(0);
  }
  50% {
    transform: translateX(-30%) translateZ(0);
  }
  100% {
    transform: translateX(0) translateZ(0);
  }
`;

const wavySvg = `
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 3200 200">
    <path fill="#003399" d="M0,100 C800,200 1600,0 2400,100 C3200,200 3200,100 3200,100 L3200,200 L0,200 Z" />
  </svg>
`;

const encodedWavySvg = encodeURIComponent(wavySvg);

const ImagePreview: React.FC<ImagePreviewProps> = ({ preview, isLoading }) => {

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', height: '100%', width: '100%', padding: '0 1rem' }}>
      {isLoading ? (
        <CircularProgress />
      ) : preview ? (
        <>
          <img src={preview} alt="Preview" style={{ maxWidth: '100%', maxHeight: '70vh' }} />
        </>
      ) : (
        <Box sx={{
          width: '100%',
          height: '100%',
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'center',
          alignItems: 'center',
          position: 'relative',
          overflow: 'hidden',
          bgcolor: '#0099ff',
          color: 'white',
          textAlign: 'center',
          p: 4,
        }}>
          <Typography variant="h4" component="h2" sx={{ 
            mb: 2, 
            fontWeight: 'bold', 
            zIndex: 1,
            textShadow: '2px 2px 4px rgba(0,0,0,0.3)',
          }}>
            No Image Selected
          </Typography>
          <Typography variant="body1" sx={{ 
            mb: 4, 
            maxWidth: '600px', 
            zIndex: 1,
            fontStyle: 'italic',
            color: '#e6f7ff',
          }}>
            Upload an image to start your watermarking journey with Watermark Wizard.
          </Typography>
          <Box sx={{
            position: 'absolute',
            left: 0,
            right: 0,
            bottom: 0,
            height: '100%',
            background: 'linear-gradient(180deg, #0099ff 0%, #0066cc 100%)',
          }}>
            <Box sx={{
              position: 'absolute',
              left: -800,
              right: -800,
              bottom: 0,
              height: '50%',
              backgroundImage: `url("data:image/svg+xml,${encodedWavySvg}")`,
              backgroundRepeat: 'repeat-x',
              backgroundSize: '3200px 200px',
              backgroundPosition: '0 bottom',
              animation: `${waveAnimation} 50s linear infinite`,
            }} />
            <Box sx={{
              position: 'absolute',
              left: -800,
              right: -800,
              bottom: '10px',
              height: '50%',
              backgroundImage: `url("data:image/svg+xml,${encodedWavySvg}")`,
              backgroundRepeat: 'repeat-x',
              backgroundSize: '3200px 200px',
              backgroundPosition: '0 bottom',
              animation: `${waveAnimation} 30s linear infinite`,
              opacity: 0.5,
            }} />
          </Box>
        </Box>
      )}
    </Box>
  );
};

export default ImagePreview;