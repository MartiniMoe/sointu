%define BPM 100

%include "../src/sointu.inc"

BEGIN_PATTERNS
    PATTERN 64, 64, 64, 64, 64,  64, 64, 64, 64, 64, 64, 64,   65, 65, 65, 65,
    PATTERN 76, 0, 0,  0, 0, 0, 0, 0, 76, 0, 0,  0, 0, 0, 0, 0,  
END_PATTERNS

BEGIN_TRACKS
    TRACK VOICES(1),0 ; a very silent voice
    TRACK VOICES(1),1 ; a loud one
END_TRACKS

BEGIN_PATCH
    BEGIN_INSTRUMENT VOICES(1) ; Instrument0
        SU_ENVELOPE MONO,ATTAC(64),DECAY(64),SUSTAIN(64),RELEASE(64),GAIN(16)
        SU_ENVELOPE MONO,ATTAC(64),DECAY(64),SUSTAIN(64),RELEASE(64),GAIN(16)
        SU_OSCILLAT MONO,TRANSPOSE(88),DETUNE(64),PHASE(0),COLOR(128),SHAPE(64),GAIN(128),FLAGS(TRISAW)
        SU_OSCILLAT MONO,TRANSPOSE(88),DETUNE(64),PHASE(0),COLOR(128),SHAPE(64),GAIN(128),FLAGS(TRISAW)
        SU_MULP     STEREO
        SU_SEND     MONO,AMOUNT(128),GLOBALPORT(2,0,receive,left) + SEND_POP
        SU_SEND     MONO,AMOUNT(128),GLOBALPORT(2,0,receive,right) + SEND_POP
    END_INSTRUMENT
    BEGIN_INSTRUMENT VOICES(1) ; Instrument0
        SU_ENVELOPE MONO,ATTAC(64),DECAY(64),SUSTAIN(64),RELEASE(64),GAIN(128)
        SU_ENVELOPE MONO,ATTAC(64),DECAY(64),SUSTAIN(64),RELEASE(64),GAIN(128)
        SU_OSCILLAT MONO,TRANSPOSE(88),DETUNE(64),PHASE(0),COLOR(128),SHAPE(64),GAIN(128),FLAGS(SINE)
        SU_OSCILLAT MONO,TRANSPOSE(88),DETUNE(64),PHASE(0),COLOR(128),SHAPE(64),GAIN(128),FLAGS(SINE)
        SU_MULP     STEREO
        SU_SEND     MONO,AMOUNT(128),GLOBALPORT(2,0,receive,left) + SEND_POP
        SU_SEND     MONO,AMOUNT(128),GLOBALPORT(2,0,receive,right) + SEND_POP
    END_INSTRUMENT
    BEGIN_INSTRUMENT VOICES(1) ; Global compressor effect
        SU_RECEIVE  STEREO
        SU_COMPRES  STEREO,ATTAC(32),RELEASE(64),INVGAIN(32),THRESHOLD(64),RATIO(96)
        SU_MULP     STEREO
        SU_OUT      STEREO, GAIN(128)
    END_INSTRUMENT
END_PATCH

%include "../src/sointu.asm"
