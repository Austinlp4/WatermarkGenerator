import React from 'react';
import { Box } from '@mui/material';

const BackgroundBlobs: React.FC = () => {
  return (
    <Box
      sx={{
        position: 'fixed',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        zIndex: -1,
        overflow: 'hidden',
        backgroundColor: 'rgba(255, 255, 255, .7)', // Add this line to set the background color to white
      }}
    >
      <svg width="100%" height="100%" xmlns="http://www.w3.org/2000/svg">
        <defs>
          <linearGradient id="blob-gradient" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" stopColor="#01204a" stopOpacity="0.3" />
            <stop offset="100%" stopColor="#01204a" stopOpacity="0.3" />
          </linearGradient>
        </defs>
        <circle cx="90%" cy="20%" r="10%" fill="url(#blob-gradient)" />
        <circle cx="30%" cy="40%" r="20%" fill="url(#blob-gradient)" />
        <circle cx="70%" cy="60%" r="15%" fill="url(#blob-gradient)" />
        <circle cx="20%" cy="80%" r="10%" fill="url(#blob-gradient)" />
        <circle cx="80%" cy="90%" r="20%" fill="url(#blob-gradient)" />
      </svg>
      <Box
        sx={{
          position: 'absolute',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          backdropFilter: 'blur(5px)',
          backgroundColor: 'rgba(255, 255, 255, 0.5)', // Adjust opacity as needed
        }}
      />
    </Box>
  );
};

export default BackgroundBlobs;
