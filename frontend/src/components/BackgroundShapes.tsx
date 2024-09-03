import React from 'react';
import { Box } from '@mui/material';

const BackgroundShapes: React.FC = () => {
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
      }}
    >
      <svg width="100%" height="100%" xmlns="http://www.w3.org/2000/svg">
        <defs>
          <linearGradient id="grad1" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" style={{ stopColor: '#8ecae6', stopOpacity: 0.5 }} />
            <stop offset="100%" style={{ stopColor: '#3a86ff', stopOpacity: 0.5 }} />
          </linearGradient>
        </defs>
        <rect width="100%" height="100%" fill="#f8f9fa" />
        <circle cx="10%" cy="10%" r="15%" fill="url(#grad1)" />
        <circle cx="90%" cy="90%" r="20%" fill="url(#grad1)" />
        <path d="M0,50 Q50,0 100,50 T200,50" stroke="url(#grad1)" strokeWidth="2" fill="none" />
      </svg>
    </Box>
  );
};

export default BackgroundShapes;
