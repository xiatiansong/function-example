#ifndef __IRZIP__H__
#define __IRZIP__H__

#define MAX_CHARARISTOR_UNIT 		12
#define MAX_ZIP_UNIT				1500
#define MAX_IR_RX_LEN				(MAX_ZIP_UNIT*2)

void irzip(char *pi8IrData, char *pi8ChararistorUnit, char *pi8ZipData);

#endif
