import { keyframes } from '@emotion/react';
import styled from '@emotion/styled';

const rotate = keyframes`
    to {
        transform: rotate(360deg);
    }
`;

const pathDash = keyframes`
    0% {
        stroke-dasharray: 1, 200;
        stroke-dashoffset: 0;
    }
    50% {
        stroke-dasharray: 89, 200;
        stroke-dashoffset: -35;
    }
    100% {
        stroke-dasharray: 89, 200;
        stroke-dashoffset: -124;
    }
`;

const pathColor = keyframes`
    0%, 100% {
        stroke: #d62d20;
    }
    40% {
        stroke: #0057e7;
    }
    66% {
        stroke: #008744;
    }
    80%, 90% {
        stroke: #ffa700;
    }
`;

const LoadingWrapper = styled.div`
    transform-origin: center center;
    width: 30px;
    height: 30px;
    margin: auto;
    animation: ${rotate} 2s linear infinite;
`;

const LoadingPath = styled.circle`
    fill: none;
    stroke-width: 3;
    stroke-linecap: round;
    animation: ${pathDash} 1.5s ease-in-out infinite,
              ${pathColor} 6s ease-in-out infinite;
`;

export const LoadingOverlay = () => {
    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-white/60">
            <LoadingWrapper>
                <svg viewBox="0 0 50 50">
                    <LoadingPath cx="25" cy="25" r="20" />
                </svg>
            </LoadingWrapper>
        </div>
    );
};

// 纯色旋转加载
const pureRotate = keyframes`
    0% {
        transform: rotate(0deg);
    }
    100% {
        transform: rotate(360deg);
    }
`;

const PureLoading = styled.div`
    width: 30px;
    height: 30px;
    border: 3px solid #FFFFFFFF;
    border-bottom-color: #84c56a;
    border-radius: 50%;
    display: inline-block;
    box-sizing: border-box;
    animation: ${pureRotate} 0.75s linear infinite;
`;

export const PureLoadingOverlay = () => {
    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-white/70">
            <PureLoading />
        </div>
    );
};
