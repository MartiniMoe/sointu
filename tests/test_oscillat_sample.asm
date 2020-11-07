%define BPM 100

%include "sointu/header.inc"

BEGIN_PATTERNS
    PATTERN 0,0,0,0,0,0,0,0
    PATTERN 72, HLD, HLD, HLD, HLD, HLD, HLD, 0
    PATTERN 64, HLD, HLD, HLD, HLD, HLD, HLD, 0
    PATTERN 60, HLD, HLD, HLD, HLD, HLD, HLD, 0
    PATTERN 40, HLD, HLD, HLD, HLD, HLD, HLD, 0
END_PATTERNS

BEGIN_TRACKS
    TRACK   VOICES(1),1,0,2,0,3,0,4,0
    TRACK   VOICES(1),0,1,0,2,0,3,0,4 ; an ordinary sine oscillator, to compare we calculate the pitch right
END_TRACKS

BEGIN_PATCH
    BEGIN_INSTRUMENT VOICES(1) ; Instrument0
        SU_ENVELOPE MONO,ATTACK(32),DECAY(32),SUSTAIN(64),RELEASE(64),GAIN(128)
        SU_ENVELOPE MONO,ATTACK(32),DECAY(32),SUSTAIN(64),RELEASE(64),GAIN(128)
        SU_OSCILLAT MONO,TRANSPOSE(68),DETUNE(64),PHASE(64),COLOR(0),SHAPE(64),GAIN(128), FLAGS(SAMPLE)
        SU_OSCILLAT MONO,TRANSPOSE(66),DETUNE(64),PHASE(64),COLOR(1),SHAPE(64),GAIN(128), FLAGS(SAMPLE)
        SU_MULP     STEREO
        SU_OUT      STEREO,GAIN(128)
    END_INSTRUMENT
    BEGIN_INSTRUMENT VOICES(1) ; Instrument1 to compare that the pitch is ok
        SU_ENVELOPE MONO,ATTACK(32),DECAY(32),SUSTAIN(64),RELEASE(64),GAIN(128)
        SU_ENVELOPE MONO,ATTACK(32),DECAY(32),SUSTAIN(64),RELEASE(64),GAIN(128)
        SU_OSCILLAT MONO,TRANSPOSE(64),DETUNE(64),PHASE(0),COLOR(128),SHAPE(64),GAIN(128), FLAGS(SINE)
        SU_OSCILLAT MONO,TRANSPOSE(64),DETUNE(64),PHASE(0),COLOR(128),SHAPE(64),GAIN(128), FLAGS(SINE)
        SU_MULP     STEREO
        SU_OUT      STEREO,GAIN(128)
    END_INSTRUMENT
END_PATCH

BEGIN_SAMPLE_OFFSETS
    SAMPLE_OFFSET START(1678611),LOOPSTART(1341),LOOPLENGTH(106) ; name VIOLN68, unitynote 56 (transpose to 4), data length 1448
    SAMPLE_OFFSET START(1680142),LOOPSTART(1483),LOOPLENGTH(95) ; name VIOLN70, unitynote 58 (transpose to 2), data length 1579
END_SAMPLE_OFFSETS

%include "sointu/footer.inc"
