import React, { useEffect, useState } from 'react';
import { Box, CircularProgress, Typography, ImageList, ImageListItem } from '@mui/material';
import { keyframes } from '@emotion/react';

interface ImagePreviewProps {
  previews?: string[];
  isLoading: boolean;
  watermarkedImages?: string[];
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

const addWatermarkToPreview = (imageUrl: string): Promise<string> => {
  console.log('addWatermarkToPreview called with:', imageUrl);
  return new Promise((resolve, reject) => {
    const img = new Image();
    img.crossOrigin = "Anonymous";
    img.onload = () => {
      console.log('Image loaded, dimensions:', img.width, 'x', img.height);
      const canvas = document.createElement('canvas');
      const ctx = canvas.getContext('2d');
      if (!ctx) {
        reject(new Error('Unable to get canvas context'));
        return;
      }
      canvas.width = img.width;
      canvas.height = img.height;
      ctx.drawImage(img, 0, 0);
      
      // Add more opaque overlay
      ctx.fillStyle = 'rgba(0, 0, 0, 0)';
      ctx.fillRect(0, 0, canvas.width, canvas.height);
      
      // Set up text style
      ctx.fillStyle = 'rgba(58, 134, 255, .6)';
      ctx.font = 'bold 130px Arial'; // Increased font size
      ctx.textAlign = 'center';
      ctx.textBaseline = 'middle';
      
      // Add drop shadow
      ctx.shadowColor = 'rgba(0, 0, 0, 0.1)';
      ctx.shadowBlur = 10;
      ctx.shadowOffsetX = 5;
      ctx.shadowOffsetY = 5;
      
      // Draw text (wrapped if necessary)
      const text = 'watermark-generator.com';
      const maxWidth = canvas.width * 0.9; // 90% of canvas width
      const lineHeight = 120;
      const words = text.split(' ');
      let line = '';
      let y = canvas.height / 2;

      for (let n = 0; n < words.length; n++) {
        const testLine = line + words[n] + ' ';
        const metrics = ctx.measureText(testLine);
        const testWidth = metrics.width;
        if (testWidth > maxWidth && n > 0) {
          ctx.fillText(line, canvas.width / 2, y);
          line = words[n] + ' ';
          y += lineHeight;
        } else {
          line = testLine;
        }
      }
      ctx.fillText(line, canvas.width / 2, y);
      
      console.log('Large watermark added successfully');
      resolve(canvas.toDataURL());
    };
    img.onerror = () => {
      console.error('Failed to load image');
      reject(new Error('Failed to load image'));
    };
    img.src = imageUrl;
  });
};

const ImagePreview: React.FC<ImagePreviewProps> = ({ previews = [], isLoading, watermarkedImages = [] }) => {
  const [localWatermarkedImages, setLocalWatermarkedImages] = useState<string[]>([]);

  useEffect(() => {
    const applyWatermarks = async () => {
      const watermarked = await Promise.all(watermarkedImages.map(addWatermarkToPreview));
      setLocalWatermarkedImages(watermarked);
    };
    applyWatermarks();
  }, [watermarkedImages]);

  return (
    <Box sx={{ 
      display: 'flex', 
      flexDirection: 'column', 
      alignItems: 'center', 
      justifyContent: 'center', 
      height: '100%', 
      width: '100%', 
      padding: '16px',
      boxSizing: 'border-box',
      overflow: 'hidden' // Prevent overflow
    }}>
      {isLoading ? (
        <CircularProgress />
      ) : previews.length > 0 ? (
        <Box sx={{ 
          width: '100%', 
          height: '100%', 
          overflow: 'auto' // Allow scrolling if needed
        }}>
          <ImageList 
            sx={{ 
              width: '100%',
              padding: 0,
              ...(previews.length === 1 
                ? { 
                    display: 'flex', 
                    justifyContent: 'center', 
                    alignItems: 'center',
                    height: '100%',
                    maxWidth: '600px',
                    margin: '0 auto'
                  } 
                : {
                    display: 'grid',
                    gridTemplateColumns: 'repeat(auto-fill, minmax(150px, 1fr))',
                    gap: '16px',
                  }
              )
            }} 
            cols={previews.length === 1 ? 1 : undefined} 
            rowHeight={previews.length === 1 ? 'auto' : undefined}
          >
            {previews.map((_, index) => (
              <ImageListItem 
                key={index} 
                sx={previews.length === 1 
                  ? { width: '100%', height: '100%', display: 'flex', justifyContent: 'center', alignItems: 'center' } 
                  : { aspectRatio: '1 / 1' }
                }
              >
                <img
                  src={localWatermarkedImages[index] || previews[index]}
                  alt={`Preview ${index + 1}`}
                  loading="lazy"
                  style={{ 
                    objectFit: previews.length === 1 ? 'contain' : 'cover', 
                    width: '100%', 
                    height: '100%',
                    maxHeight: previews.length === 1 ? '100%' : 'none',
                  }}
                />
              </ImageListItem>
            ))}
          </ImageList>
        </Box>
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