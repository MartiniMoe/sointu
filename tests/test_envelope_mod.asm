%define BPM 100

%include "sointu/header.inc"

BEGIN_PATTERNS
    PATTERN 64, HLD, HLD, HLD, HLD, HLD, HLD, HLD,HLD, HLD, HLD, 0, 0, 0, 0, 0
END_PATTERNS

BEGIN_TRACKS
    TRACK VOICES(1),0
END_TRACKS

BEGIN_PATCH
    BEGIN_INSTRUMENT VOICES(1) ; Instrument0
        SU_ENVELOPE MONO,ATTACK(80),DECAY(80),SUSTAIN(64),RELEASE(80),GAIN(128)
        SU_ENVELOPE MONO,ATTACK(80),DECAY(80),SUSTAIN(64),RELEASE(80),GAIN(128)
        SU_OSCILLAT MONO,TRANSPOSE(120),DETUNE(64),PHASE(0),COLOR(128),SHAPE(96),GAIN(128),FLAGS(SINE+LFO)
        SU_SEND MONO,AMOUNT(68),UNIT(0),PORT(0),FLAGS(NONE)
        SU_SEND MONO,AMOUNT(68),UNIT(0),PORT(1),FLAGS(NONE)
        ; Sustain modulation seems not to be implemented
        SU_SEND MONO,AMOUNT(68),UNIT(0),PORT(3),FLAGS(NONE)
        SU_SEND MONO,AMOUNT(68),UNIT(1),PORT(4),FLAGS(SEND_POP)
        SU_OUT  STEREO,GAIN(110)
    END_INSTRUMENT
END_PATCH

%include "sointu/footer.inc"
